package logger

import (
	"errors"
	"log/slog"
	"os"
)

const (
	envLocal  = "local"
	envDocker = "docker"
)

var ErrInvalidEnv = errors.New("invalid env type")

func New(env string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDocker:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return nil, ErrInvalidEnv
	}

	return logger, nil
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
