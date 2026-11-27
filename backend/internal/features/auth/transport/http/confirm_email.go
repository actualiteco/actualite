package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

// ConfirmEmail godoc
// @Summary Подтверждение почты
// @Description Принимает токен для подтверждения почты (должен парситься из URL ссылки, которая присылается на почту при запросе на подтверждение)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ConfirmEmailDTORequest true "Тело запроса"
// @Success 204 "Почта успешно подтверждена"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 401 {object} domain.CustomError "Ошибка авторизации"
// @Failure 404 {object} domain.CustomError "Пользователь не найден"
// @Failure 409 {object} domain.CustomError "Не удалось подтвердить почту"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/email-confirmation/confirm [post]
// @Security BearerAuth
func (h *AuthHTTPHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req ConfirmEmailDTORequest
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		respWriter.MapError(err)
		return
	}

	userID, err := middleware.GetUserUUIDFromCtx(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	if err := h.AuthService.ConfirmEmail(ctx, userID, req.Token); err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.NoContent()
}
