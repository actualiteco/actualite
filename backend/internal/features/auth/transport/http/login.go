package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"

	"go.uber.org/zap"
)

// LoginUser godoc
// @Summary Авторизация по email и паролю
// @Description Авторизует пользователя по email и паролю, возвращает access и refresh токены с временем жизни access token в секундах (10 минут)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginDTO true "Тело запроса"
// @Success 200 {object} TokensDTOResponse "Успешная авторизация"
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

	ip, err := httprequest.GetIP(r)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	accessToken, refreshToken, expirationTime, err := h.AuthService.Login(ctx, req.Email, req.Password, r.UserAgent(), ip)
	if err != nil {
		log.Error("login user", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(toTokensDTOResponse(accessToken, refreshToken, expirationTime), http.StatusOK)
}
