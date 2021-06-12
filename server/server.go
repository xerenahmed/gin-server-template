package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	cache2 "server-app/server/manager/cache"
	db2 "server-app/server/manager/db"
	router2 "server-app/server/manager/router"
	settings2 "server-app/server/settings"
	"time"
)

type Server struct {
	db       *db2.Manager
	cache    *cache2.Manager
	router   *router2.Manager
	settings *settings2.Settings

	CloseReceiver chan struct{}
}

func NewServer() (*Server, error) {
	var (
		settings *settings2.Settings
		err      error
	)

	settings, err = settings2.Load()
	if err != nil {
		return nil, fmt.Errorf("failed loading settings: %v", err)
	}

	db := db2.NewDatabaseManager()
	cache := cache2.NewCacheManager()

	return &Server{
		db:        db,
		cache:     cache,
		settings:  settings,
		router:    router2.NewRouterManager(db, cache, settings),
		CloseReceiver: make(chan struct{}, 1),
	}, nil
}

func (s *Server) Load() {
	var err error

	go func() {
		logrus.Info("connecting to redis")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()

		_, err := s.cache.Ping(ctx)
		if err != nil {
			logrus.WithError(err).Warning("failed ping redis")
		}
	}()

	logrus.Info("connecting to mongodb")
	err = s.db.Load()
	if err != nil {
		logrus.WithError(err).Fatal("failed load mongodb.")
	}

	logrus.Info("router handlers loading")
	err = s.router.Load()
	if err != nil {
		logrus.WithError(err).Fatal("failed loading router handlers")
	}

	logrus.Info("server ready")
}

func (s *Server) Run() {
	errReceiver := make(chan error, 1)
	go s.router.Run(errReceiver)

	select {
	case err := <-errReceiver:
		logrus.WithError(err).Error("an error occurred")
		s.Close()
	}
}

func (s *Server) Close() {
	s.shutdown()
	close(s.CloseReceiver)
}

func (s *Server) shutdown() {
	err := s.router.Shutdown()
	if err != nil {
		logrus.WithError(err).Error("failed stopping router handlers")
	}
}

func (s Server) Settings() settings2.Settings {
	return *s.settings
}