package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server-app/middleware"
	"server-app/router"
	"server-app/router/auth"
	"server-app/server/manager/cache"
	"server-app/server/manager/db"
	"server-app/server/settings"
	"time"
)

type Manager struct {
	handler *gin.Engine
	server  *http.Server

	db       *db.Manager
	cache    *cache.Manager
	settings *settings.Settings
}

func NewRouterManager(db *db.Manager, c *cache.Manager, s *settings.Settings) *Manager {
	handler := gin.New()
	server := &http.Server{
		Addr:    ":5000",
		Handler: handler,
	}

	return &Manager{handler, server, db, c, s}
}

func (m *Manager) Load() error {
	m.handler.Use(middleware.Cors)
	m.handler.GET("/", router.Index)

	m.loadAdmin()
	m.loadAuth()
	m.loadGeneral()

	return nil
}

func (m *Manager) Run(errReceiver chan error) {
	if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errReceiver <- err
	}
}

func (m *Manager) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("yönlendiriciler düzgün bir şekilde kapanmadı: %v", err)
	}

	return nil
}

func (m *Manager) loadGeneral() {
	_ = m.handler.Group("")
}

func (m *Manager) loadAdmin() {
	mod := m.handler.Group("/admin")
	mod.Use(middleware.AuthRequired)
}

func (m *Manager) loadCdn() {
	_ = m.handler.Group("/cdn")
}

func (m *Manager) loadAuth() {
	authRouter := m.handler.Group("/auth")
	authRouter.POST("/login", auth.Login(m.db))
}
