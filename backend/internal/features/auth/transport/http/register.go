package auth_transport

import (
	"encoding/json"
	"net/http"
	"time"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"

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

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя и возвращает его данные
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterDTO true "тело запроса"
// @Success 201 {object} RegisterDTOResponse "Пользователь создан"
// @Failure 400 {object} domain.CustomError "Неверный запрос или email уже занят"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *AuthHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("error reguster decoding", zap.Error(err))

		respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
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
