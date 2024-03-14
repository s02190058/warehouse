package http

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/s02190058/warehouse/internal/transport/http/response"
)

func ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.OK())
}
