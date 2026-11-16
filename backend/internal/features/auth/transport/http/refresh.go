package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" example:"old_refresh_token" validate:"required"`
}

// Refresh godoc
// @Summary Обновление refresh-токена
// @Description Обновляет refresh-токен, возвращает access и refresh токены с временем жизни access token в секундах (10 минут). Используется, когда прежний AccessToken протух
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshDTO true "Тело запроса"
// @Success 200 {object} TokensDTOResponse "Успешное обновление токена"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 401 {object} domain.CustomError "Refresh-токен протух или не валиден"
// @Failure 404 {object} domain.CustomError "Сессия или пользователь-обладатель сессии не найден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/refresh [post]
func (h *AuthHTTPHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req RefreshDTO
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		respWriter.MapError(err)
		return
	}

	accessToken, refreshToken, expirationTime, err := h.AuthService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(toTokensDTOResponse(accessToken, refreshToken, expirationTime), http.StatusOK)
}
