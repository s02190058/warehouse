package postgres

import "time"

type Option func(options *Options)

func WithConnAttempts(attempts int) Option {
	return func(o *Options) {
		o.connAttempts = attempts
	}
}

func WithConnTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.connTimeout = timeout
	}
}

func WithMaxConns(conns int32) Option {
	return func(o *Options) {
		o.maxConns = conns
	}
}
