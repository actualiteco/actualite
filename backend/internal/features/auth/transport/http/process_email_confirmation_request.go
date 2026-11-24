package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

// ProcessEmailConfirmationRequest godoc
// @Summary Обработка запроса на подтверждение email
// @Description Создает токен, его хэш сохраняется в БД, а на почту приходит письмо с ссылкой на страницу фронтенда (/confirm-email?token=...) для подтверждения email
// @Tags Auth
// @Accept json
// @Produce json
// @Success 202 "Письмо с ссылкой для верификации email отправлено на почту"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 401 {object} domain.CustomError "Пользователь не авторизован"
// @Failure 404 {object} domain.CustomError "Пользователь не найден"
// @Failure 409 {object} domain.CustomError "Почта уже подтверждена"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/email-confirmation/request [post]
// @Security BearerAuth
func (h *AuthHTTPHandler) ProcessEmailConfirmationRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	if err := h.AuthService.ProcessEmailConfirmationRequest(ctx, userID); err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.Accepted()
}
