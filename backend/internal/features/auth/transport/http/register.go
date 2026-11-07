package auth_transport

import (
	"net/http"
	"time"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email" example:"m4rkek@gmail.com" extensions:"x-order=0"`
	Username string `json:"username" example:"m4rkek" extensions:"x-order=1"`
	Password string `json:"password" validate:"required,min=8,max=64,ascii" example:"supersecretpass" extensions:"x-order=2"`
	Role     string `json:"role" validate:"required" example:"admin" extensions:"x-order=3"`
}

func (r RegisterDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type RegisterDTOResponse struct {
	ID        string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	Username  string    `json:"username" example:"m4rkek" extensions:"x-order=1"`
	Email     string    `json:"email" example:"m4rkek@gmail.com" extensions:"x-order=2"`
	Role      string    `json:"role" example:"user" extensions:"x-order=3"`
	CreatedAt time.Time `json:"createdAt" example:"2026-03-25T12:00:41.267Z" extensions:"x-order=4"`
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

	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req RegisterDTO
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		log.Debug("decode and validate register request")
		respWriter.MapError(err)
		return
	}

	user, err := h.AuthService.RegisterUser(
		ctx,
		domain.User{
			Email:    req.Email,
			Username: req.Username,
			Role:     req.Role,
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
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	respWriter.JSONResponse(resp, http.StatusCreated)
}
