package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	w   http.ResponseWriter
}

func NewResponseHandler(
	l *logger.Logger,
	w http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: l,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(resp any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(resp); err != nil {
		h.log.Error("write json response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) errorResponse(err domain.CustomError, statusCode int) {
	h.JSONResponse(map[string]domain.CustomError{"error": err}, statusCode)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(domain.ErrInternal, http.StatusInternalServerError)
}

func (h *HTTPResponseHandler) MapError(err error) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
		customErr  domain.CustomError
	)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Info
		customErr = domain.ErrNotFound
	case errors.Is(err, domain.ErrInvalidRequest):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Error
		customErr = domain.ErrInvalidRequest
	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn
		customErr = domain.ErrUnauthorized
	case errors.Is(err, domain.ErrEmailAlreadyUsed):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
		customErr = domain.ErrEmailAlreadyUsed
	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn
		customErr = domain.ErrForbidden
	case errors.Is(err, domain.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
		customErr = domain.ErrConflict
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
		customErr = domain.ErrInternal
	}

	logFunc(err.Error(), zap.Error(err))

	h.errorResponse(customErr, statusCode)
}
