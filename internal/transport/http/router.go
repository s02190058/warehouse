package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	// register swagger docs.
	_ "github.com/s02190058/warehouse/docs"
)

//	@title			warehouse App
//	@version		0.1.o
//	@description	Warehouse management platform.

//	@contact.name	Bakanov Artem
//	@contact.url	https://t.me/s02190058
//	@contact.email	sklirian@mail.ru

// @BasePath	/

func NewRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Get("/ping", ping)

	return router
}
