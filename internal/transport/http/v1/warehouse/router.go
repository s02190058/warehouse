package warehouse

import (
	"log/slog"

	"github.com/go-chi/chi"
)

func NewRouter(
	logger *slog.Logger,
	service Service,
) chi.Router {
	h := &handler{
		logger:  logger,
		service: service,
	}

	router := chi.NewRouter()

	router.Get("/{id}", h.remains)
	router.Post("/{id}:reserve", h.reserve)
	router.Post("/{id}:release", h.release)

	return router
}
