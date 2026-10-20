package auth_transport

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"
	"github.com/tryingmyb3st/PolyTweet/internal/core/logger"
	"github.com/tryingmyb3st/PolyTweet/internal/core/transport/response"

	"go.uber.org/zap"
)

type RegisterDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterDTOResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *AuthHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("error reguster decoding", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.RegisterUser(
		ctx,
		domain.User{
			Email: req.Email,
			Role:  req.Role,
		},
		req.Password,
	)

	if err != nil {
		log.Error("get user after register", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := RegisterDTOResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	respWriter.JSONResponse(resp, http.StatusCreated)
}
