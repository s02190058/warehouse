package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	// register swagger docs.
	_ "github.com/s02190058/warehouse/docs"
	v1 "github.com/s02190058/warehouse/internal/transport/http/v1"
	"github.com/s02190058/warehouse/internal/transport/http/v1/warehouse"
)

//	@title			warehouse App
//	@version		0.1.o
//	@description	Warehouse management platform.

//	@contact.name	Bakanov Artem
//	@contact.url	https://t.me/s02190058
//	@contact.email	sklirian@mail.ru

// @BasePath	/

func NewRouter(
	logger *slog.Logger,
	warehouseService warehouse.Service,
) chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Get("/ping", ping)

	api := chi.NewRouter()

	api.Mount("/v1", v1.NewRouter(logger, warehouseService))

	router.Mount("/api", api)

	return router
}
