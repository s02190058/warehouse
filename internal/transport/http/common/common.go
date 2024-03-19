package common

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	operationAttrKey = "operation"
	requestIDAttrKey = "request_id"
)

func LoggerWithHandler(logger *slog.Logger, op string, r *http.Request) *slog.Logger {
	return logger.With(
		slog.String(operationAttrKey, op),
		slog.String(requestIDAttrKey, middleware.GetReqID(r.Context())),
	)
}
