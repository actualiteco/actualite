package auth_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

// ConfirmPasswordReset godoc
// @Summary Подтверждение смены пароля
// @Description Принимает токен для смены пароля и новый пароль, если все успешно - возвращается 204
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ConfirmPasswordResetDTORequest true "Тело запроса"
// @Success 204 "Пароль успешно подтвержден"
// @Failure 400 {object} domain.CustomError "Невалидный запрос"
// @Failure 409 {object} domain.CustomError "Не удалось подтвердить смену пароля"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/password-reset/confirm [post]
func (h *AuthHTTPHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req ConfirmPasswordResetDTORequest
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		respWriter.MapError(err)
		return
	}

	if err := h.AuthService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.NoContent()
}
