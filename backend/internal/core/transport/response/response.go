package response

import (
	"encoding/json"
	"errors"
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

func (h *HTTPResponseHandler) ErrorResponse(err domain.CustomError, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	resp := map[string]domain.CustomError{
		"error": err,
	}

	if err := json.NewEncoder(h.w).Encode(resp); err != nil {
		h.log.Error("write error response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(http.StatusInternalServerError)

	response := map[string]domain.CustomError{
		"error": domain.ErrInternal,
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("write panic response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) MapError(err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		h.ErrorResponse(domain.ErrInvalidRequest, http.StatusNotFound)
		return
	case errors.Is(err, domain.ErrInvalidRequest):
		h.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
		return
	case errors.Is(err, domain.ErrUnauthorized):
		h.ErrorResponse(domain.ErrUnauthorized, http.StatusUnauthorized)
		return
	case errors.Is(err, domain.ErrEmailAlreadyUsed):
		h.ErrorResponse(domain.ErrEmailAlreadyUsed, http.StatusConflict)
	case errors.Is(err, domain.ErrForbidden):
		h.ErrorResponse(domain.ErrForbidden, http.StatusForbidden)
		return
	default:
		h.ErrorResponse(domain.ErrInternal, http.StatusInternalServerError)
	}
}
