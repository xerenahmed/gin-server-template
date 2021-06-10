package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type GeneralDataBase struct {
	db *mongo.Database
}

func (m GeneralDataBase) Some() *mongo.Collection {
	return m.db.Collection("some")
}

func (m GeneralDataBase) CreateCollections() error {
	collections := []*mongo.Collection{
		m.Some(),
	}

	for _, col := range collections {
		if err := m.db.CreateCollection(context.Background(), col.Name());
			err != nil && !strings.HasPrefix(err.Error(), "(NamespaceExists) Collection already exists") {
			return fmt.Errorf("can't create collection \"%s\": %v", col.Name(), err)
		}
	}

	return nil
}

func (m GeneralDataBase) CreateIndexes() error {
	return nil
}

func (m GeneralDataBase) CreateValidators() error {
	return nil
}

func (m GeneralDataBase) Name() string {
	return m.db.Name()
}