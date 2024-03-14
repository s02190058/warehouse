package app

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/s02190058/warehouse/internal/config"
	"github.com/s02190058/warehouse/internal/transport/http"
	httpserver "github.com/s02190058/warehouse/pkg/http/server"
	"github.com/s02190058/warehouse/pkg/slogger"
)

func Run(cfg *config.Config) {
	logger, err := slogger.New(cfg.Env)
	if err != nil {
		log.Fatalf("cannot get logger: %v", err)
	}

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
