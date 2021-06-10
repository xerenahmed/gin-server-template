package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sync/errgroup"
	"os"
	"time"
)

type DatabaseManager interface {
	CreateCollections() error
	CreateValidators() error
	CreateIndexes() error
	Name() string
}

type Manager struct {
	general  *GeneralDataBase
	logs     *LogDataBase

	client *mongo.Client
}

func NewDatabaseManager() *Manager {
	return &Manager{}
}

func (m *Manager) Load() error {
	var err error

	for i := 0; i < 5; i++ {
		m.client, err = createClient()
		if err != nil {
			if i == 2 {
				return err
			}
		} else {
			break
		}
		time.Sleep(time.Second * 3)
	}

	m.general = &GeneralDataBase{db: m.client.Database("general")}
	m.logs = &LogDataBase{db: m.client.Database("logs")}

	dbManagers := []DatabaseManager{
	 	m.general, m.logs,
	}
	errs := errgroup.Group{}
	for _, dbManager := range dbManagers {
		errs.Go(dbManager.CreateCollections)
	}
	if err := errs.Wait(); err != nil {
		return err
	}

	for _, dbManager := range dbManagers {
		errs.Go(dbManager.CreateIndexes)
	}
	if err := errs.Wait(); err != nil {
		return err
	}

	for _, dbManager := range dbManagers {
		errs.Go(dbManager.CreateValidators)
	}
	if err := errs.Wait(); err != nil {
		return err
	}

	return nil
}

func createClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_ADDR")))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}

func (m *Manager) General() *GeneralDataBase {
	return m.general
}

func (m *Manager) Logs() *LogDataBase {
	return m.logs
}

func (m *Manager) Client() *mongo.Client {
	return m.client
}
