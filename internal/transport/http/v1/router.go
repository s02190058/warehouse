package v1

import (
	"log/slog"

	"github.com/go-chi/chi"

	"github.com/s02190058/warehouse/internal/transport/http/v1/warehouse"
)

func NewRouter(
	logger *slog.Logger,
	warehouseService warehouse.Service,
) chi.Router {
	router := chi.NewRouter()

	router.Mount("/warehouses", warehouse.NewRouter(logger, warehouseService))

	return router
}
