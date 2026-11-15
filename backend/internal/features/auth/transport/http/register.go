package auth_transport

import (
	"net/http"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func (r RegisterDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя и возвращает его данные
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterDTO true "тело запроса"
// @Success 201 {object} RegisterDTOResponse "Пользователь создан"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 409 {object} domain.CustomError "Email или username уже заняты"
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
