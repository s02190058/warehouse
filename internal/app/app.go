package app

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/s02190058/warehouse/internal/config"
	"github.com/s02190058/warehouse/internal/transport/http"
	"github.com/s02190058/warehouse/pkg/db/postgres"
	httpserver "github.com/s02190058/warehouse/pkg/http/server"
	"github.com/s02190058/warehouse/pkg/slogger"
)

func Run(cfg *config.Config) {
	logger, err := slogger.New(cfg.Env)
	if err != nil {
		log.Fatalf("cannot get logger: %v", err)
	}

	db, err := postgres.New(logger, postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DB:       cfg.Postgres.DB,
	})
	if err != nil {
		logger.Error("cannot establish connection to postgresql server", slogger.Err(err))
		return
	}

	logger.Info("established connection to postgresql server")

	_ = db

	router := http.NewRouter()

	port := cfg.HTTP.Port
	server := httpserver.New(router, port)

	server.Start()
	logger.Info(
		"started http server",
		slog.String("port", port),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Info("service interrupt", slogger.Err(ctx.Err()))
	case err = <-server.Notify():
		logger.Error("error occurred during http server running", slogger.Err(err))
	}

	if err = server.Shutdown(); err != nil {
		logger.Error("error occurred during server shutdown", slogger.Err(err))
	}
}
