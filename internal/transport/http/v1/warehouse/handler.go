package warehouse

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	entity "github.com/s02190058/warehouse/internal/entity/warehouse"
	service "github.com/s02190058/warehouse/internal/service/warehouse"
	"github.com/s02190058/warehouse/internal/transport/http/common"
	"github.com/s02190058/warehouse/internal/transport/http/response"
	"github.com/s02190058/warehouse/pkg/slogger"
)

type Service interface {
	Remains(ctx context.Context, id int) (remains []entity.ProductRemains, err error)
	Reserve(ctx context.Context, id int, productCodes []string) (reservedCodes []string, err error)
	Release(ctx context.Context, id int, productCodes []string) (releasedCodes []string, err error)
}

type handler struct {
	logger *slog.Logger

	service Service
}

//	@Summary		OK status
//	@Description	Number of remaining products.
//	@Tags			warehouse
//	@Produce		json
//	@Param			id	path	int	true	"warehouse id"
//	@Success		200	{array}	warehouse.ProductRemains
//	@Router			/api/v1/warehouses/{id} [get]
//
// remains shows number of remaining products.
func (h *handler) remains(w http.ResponseWriter, r *http.Request) {
	const op = "warehouse.remains"

	logger := common.LoggerWithHandler(h.logger, op, r)

	idRaw := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		logger.Info("failed to convert id", slogger.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(ErrBadURLParam))

		return
	}

	remains, err := h.service.Remains(r.Context(), id)
	if err != nil {
		render.Status(r, error2StatusCode(err))
		render.JSON(w, r, response.Error(err))

		return
	}

	render.JSON(w, r, remains)
}

//	@Summary		OK status
//	@Description	Reserves products with the specified codes.
//	@Tags			warehouse
//	@Accept			json
//	@Produce		json
//	@Param			id				path	int			true	"warehouse id"
//	@Param			reservedCodes	body	[]string	true	"product	codes	to	be	reserved"
//	@Success		200				{array}	string
//	@Router			/api/v1/warehouses/{id}:reserve [post]
//
// reserve reserves products with the specified codes.
func (h *handler) reserve(w http.ResponseWriter, r *http.Request) {
	const op = "warehouse.reserve"

	logger := common.LoggerWithHandler(h.logger, op, r)

	var productCodes []string
	if err := render.DecodeJSON(r.Body, &productCodes); err != nil {
		logger.Info("failed to parse request body", slogger.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(ErrBadJSONBody))

		return
	}

	if err := r.Body.Close(); err != nil {
		logger.Error("failed to close request body", slogger.Err(err))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(ErrInternal))

		return
	}

	idRaw := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		logger.Info("failed to convert id", slogger.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(ErrBadURLParam))

		return
	}

	reservedCodes, err := h.service.Reserve(r.Context(), id, productCodes)
	if err != nil {
		render.Status(r, error2StatusCode(err))
		render.JSON(w, r, response.Error(err))

		return
	}

	render.JSON(w, r, reservedCodes)
}

//	@Summary		OK status
//	@Description	Release products with the specified codes.
//	@Tags			warehouse
//	@Accept			json
//	@Produce		json
//	@Param			id				path	int			true	"warehouse id"
//	@Param			releasedCodes	body	[]string	true	"product	codes	to	be	released"
//	@Success		200				{array}	string
//	@Router			/api/v1/warehouses/{id}:release [post]
//
// release releases products with the specified codes.
func (h *handler) release(w http.ResponseWriter, r *http.Request) {
	const op = "warehouse.release"

	logger := common.LoggerWithHandler(h.logger, op, r)

	var productCodes []string
	if err := render.DecodeJSON(r.Body, &productCodes); err != nil {
		logger.Info("failed to parse request body", slogger.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(ErrBadJSONBody))

		return
	}

	if err := r.Body.Close(); err != nil {
		logger.Error("failed to close request body", slogger.Err(err))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error(ErrInternal))

		return
	}

	idRaw := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		logger.Info("failed to convert id", slogger.Err(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(ErrBadURLParam))

		return
	}

	releasedCodes, err := h.service.Release(r.Context(), id, productCodes)
	if err != nil {
		render.Status(r, error2StatusCode(err))
		render.JSON(w, r, response.Error(err))

		return
	}

	render.JSON(w, r, releasedCodes)
}

func error2StatusCode(err error) (code int) {
	switch {
	case errors.Is(err, service.ErrWarehouseNotFound):
		code = http.StatusNotFound
	case errors.Is(err, service.ErrWarehouseNotAvailable):
		code = http.StatusForbidden
	case errors.Is(err, service.ErrProductIsFullyReserved) ||
		errors.Is(err, service.ErrProductHasNoReserves):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	return
}
