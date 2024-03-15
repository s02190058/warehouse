package http

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/s02190058/warehouse/internal/transport/http/response"
)

//	@Summary		OK status
//	@Description	Shows that service is available.
//	@Tags			healthcheck
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/ping [get].
func ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.OK())
}
