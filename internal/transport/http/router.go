package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
	)

	router.Get("/ping", ping)

	return router
}
