package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

// ProcessResetPasswordRequest godoc
// @Summary Обработка запроса на смену пароля
// @Description Создает password_reset_token с временем жизни, указанным в .env, а также асинхронно присылает письмо с ссылкой на страницу фронтенда для смены пароля (если почта есть и верифицирована. если нет - статус-код ответа также 202, но письмо не будет отправлено)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordDTORequest true "Тело запроса"
// @Success 202 "Письмо с ссылкой для смены пароля отправлено на почту, если почта верифицирована"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/password-reset/request [post]
func (h *AuthHTTPHandler) ProcessResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req ResetPasswordDTORequest
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		respWriter.MapError(err)
		return
	}

	if err := h.AuthService.ProcessResetPasswordRequest(ctx, req.Email); err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.Accepted()
}
