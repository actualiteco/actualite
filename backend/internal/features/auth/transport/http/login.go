package auth_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type LoginDTO struct {
	Email    string `json:"email" example:"test.user@email.com" validate:"required"`
	Password string `json:"password" example:"somepass" validate:"required"`
}

type LoginDTOResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// LoginUser godoc
// @Summary Авторизация по email и паролю
// @Description Авторизует пользователя по email и паролю, возвращает JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginDTO true "тело запроса"
// @Success 200 {object} LoginDTOResponse "Успешная авторизация"
// @Failure 401 {object} domain.CustomError "Неверные учётные данные"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /login [post]
func (h *AuthHTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("error register decoding", zap.Error(err))

		respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	jwt, user, err := h.AuthService.LoginUser(ctx, req.Email, req.Password)

	if err != nil {
		log.Error("get user", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := LoginDTOResponse{
		Token:    *jwt,
		Username: user.Username,
		UserID:   user.ID,
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}
