package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"

	"go.uber.org/zap"
)

type LoginDTO struct {
	Email    string `json:"email" example:"test.user@email.com" validate:"required"`
	Password string `json:"password" example:"hashpassword" validate:"required,min=8,max=64,ascii"`
}

type LoginDTOResponse struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTM2NjUxNzcsImlhdCI6MTc5MzYyMTk3NywiVXNlcklkIjoiODc5NGU1ODktMGRkYi00M2NlLTlmOTItMTZmYWFmY2Y0ZWU0IiwiUm9sZSI6InVzZXIifQ.iv_JT1Fke-QaK2nx2FEF_pkDroyJP-_xK81egOAV7h4"`
	UserID   string `json:"user_id" example:"8794e589-0ddb-43ce-9f92-16faafcf4ee4"`
	Username string `json:"username" example:"username"`
}

// LoginUser godoc
// @Summary Авторизация по email и паролю
// @Description Авторизует пользователя по email и паролю, возвращает JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginDTO true "тело запроса"
// @Success 200 {object} LoginDTOResponse "Успешная авторизация"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учётные данные"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /login [post]
func (h *AuthHTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req LoginDTO
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		log.Debug("decode and validate login request", zap.Error(err))
		respWriter.MapError(err)
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
