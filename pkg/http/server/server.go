package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second

	_defaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	httpServer      *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, port string, opts ...Option) *Server {
	server := newDefaultServer(handler, port)

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func newDefaultServer(handler http.Handler, port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         net.JoinHostPort("", port),
			Handler:      handler,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.httpServer.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
