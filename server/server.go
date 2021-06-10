package server

import (
	"context"
	"fmt"
	cache2 "server-app/server/manager/cache"
	db2 "server-app/server/manager/db"
	router2 "server-app/server/manager/router"
	settings2 "server-app/server/settings"
	"github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	db       *db2.Manager
	cache    *cache2.Manager
	router   *router2.Manager
	settings *settings2.Settings

	CloseRecv chan struct{}
}

func NewServer() (*Server, error) {
	var (
		settings *settings2.Settings
		err      error
	)

	settings, err = settings2.Load()
	if err != nil {
		return nil, fmt.Errorf("ayarlar yüklenemedi: %v", err)
	}

	db := db2.NewDatabaseManager()
	cache := cache2.NewCacheManager()

	return &Server{
		db:        db,
		cache:     cache,
		settings:  settings,
		router:    router2.NewRouterManager(db, cache, settings),
		CloseRecv: make(chan struct{}, 1),
	}, nil
}

func (s *Server) Load() {
	var err error

	logrus.Info("Veritabanı yükleniyor.")
	err = s.db.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Veritabanı yüklenemedi.")
	}

	go func() {
		logrus.Info("Cache servisine bağlanılıyor.")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()
		_, err = s.cache.Ping(ctx)
		if err != nil {
			logrus.WithError(err).Warning("Cache servisi ile bağlantı kurulamadı.")
		}
	}()

	logrus.Info("Yönlendiriciler başlatılıyor.")
	err = s.router.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Yönlendiriciler başlatılamadı.")
	}

	logrus.Info("Sunucu hazır!")
}

func (s *Server) Run() {
	errorChan := make(chan error, 1)
	go s.router.Run(errorChan)

	select {
	case err := <-errorChan:
		logrus.WithError(err).Error("Bir hata oluştu.")
		s.Close()
	}
}

func (s *Server) Close() {
	close(s.CloseRecv)
}

func (s *Server) Shutdown() {
	err := s.router.Shutdown()
	if err != nil {
		logrus.WithError(err).Error("Yönlendiriciler kapatılırken bir hata oluştu.")
	}
}

func (s Server) Settings() settings2.Settings {
	return *s.settings
}