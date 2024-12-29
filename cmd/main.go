package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rusystem/cache"
	"github.com/rusystem/crm-api/internal/config"
	"github.com/rusystem/crm-api/internal/repository"
	http_server "github.com/rusystem/crm-api/internal/server/http"
	"github.com/rusystem/crm-api/internal/service"
	http_handler "github.com/rusystem/crm-api/internal/transport/http"
	"github.com/rusystem/crm-api/pkg/auth"
	"github.com/rusystem/crm-api/pkg/database"
	"github.com/rusystem/crm-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// init logger
func init() {
	logger.ZapLoggerInit()
}

// @title Web api gateway API
// @version 1.0
// @description API gateway

// @contact.name ru.system.ru@gmail.com
// @contact.email ru.system.ru@gmail.com

// @host 91.243.71.100:8080
// @BasePath /api/web-api-gateway/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// init configs
	cfg, err := config.New(false)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize config, err: %v", err))
	}

	// init in-memory cache
	memCache := cache.New()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize cache, err: %v", err))
	}

	// init token manager
	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to initialize token manager, err: %v", err))
	}

	// init postgres connection
	pc, err := database.NewPostgresConnection(database.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	defer func(pc *sql.DB) {
		if err = pc.Close(); err != nil {
			logger.Error(fmt.Sprintf("postgres: failed to close connection, err: %v", err.Error()))
		}
	}(pc)

	// init dep-s
	repo := repository.New(cfg, memCache, pc)
	srv := service.New(service.Config{
		Config:       cfg,
		Repo:         repo,
		TokenManager: tokenManager,
	})
	hh := http_handler.NewHandler(srv, tokenManager, cfg)

	// HTTP Server
	server := http_server.New(cfg, hh.Init())

	go func() {
		if err = server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(fmt.Sprintf("error occurred while running http server: %s", err))
		}
	}()

	logger.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = server.Stop(ctx); err != nil {
		logger.Error(fmt.Sprintf("failed to stop server: %v", err))
	}
}
