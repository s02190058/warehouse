package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 5 * time.Second
	_defaultMaxConns     = 4
)

type Options struct {
	connAttempts int
	connTimeout  time.Duration
	maxConns     int32
}

type Database struct {
	Pool *pgxpool.Pool
}

func New(logger *slog.Logger, cfg Config, opts ...Option) (*Database, error) {
	o := newDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	url := cfg.URL()

	logger.Debug(
		"trying to connect to postgresql server",
		slog.String("url", url),
		slog.Int("attempts", o.connAttempts),
		slog.Duration("timeout", o.connTimeout),
	)

	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = o.maxConns

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	for {
		o.connAttempts--
		if err = pool.Ping(context.Background()); err == nil || o.connAttempts == 0 {
			break
		}

		logger.Debug(
			"unable to connect to postgresql server",
			slog.Int("attempts left", o.connAttempts),
		)

		time.Sleep(o.connTimeout)
	}
	if err != nil {
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}

func newDefaultOptions() *Options {
	return &Options{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
		maxConns:     _defaultMaxConns,
	}
}
