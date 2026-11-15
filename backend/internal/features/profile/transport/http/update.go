package profile_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"

	"go.uber.org/zap"
)

// Update godoc
// @Summary Обновить профиль пользователя
// @Description Обновить профиль пользователя (display_name, username, status)
// @Tags Profile
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Информация профиля"
// @Success 200 {object} ProfileResponse "Успешный запрос"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 403 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 409 {object} domain.CustomError "Пользователь с таким юзернеймом уже существует"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile [patch]
// @Security BearerAuth
func (h *ProfileHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Error("get userID from context", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	var req UpdateProfileRequest
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		log.Debug("decode and validate update profile request")
		respWriter.MapError(err)
		return
	}

	user, err := h.profileService.Update(ctx, userID, req.DisplayName, req.Username, req.Status)
	if err != nil {
		log.Error("failed to update profile", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(toResponse(user), http.StatusOK)
}
