package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"
	"github.com/tryingmyb3st/PolyTweet/internal/core/logger"
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
		"error": domain.INTERNAL_ERROR,
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("write panic response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) MapError(err error) {
	switch {
	case errors.Is(err, domain.NOT_FOUND):
		h.ErrorResponse(domain.INVALID_REQUEST, http.StatusNotFound)
		return
	case errors.Is(err, domain.INVALID_REQUEST):
		h.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	case errors.Is(err, domain.UNAUTHORIZED):
		h.ErrorResponse(domain.UNAUTHORIZED, http.StatusUnauthorized)
		return
	case errors.Is(err, domain.EMAIL_ALREADY_USED):
		h.ErrorResponse(domain.EMAIL_ALREADY_USED, http.StatusConflict)
	case errors.Is(err, domain.FORBIDDEN):
		h.ErrorResponse(domain.FORBIDDEN, http.StatusForbidden)
		return
	default:
		h.ErrorResponse(domain.INTERNAL_ERROR, http.StatusInternalServerError)
	}
}
