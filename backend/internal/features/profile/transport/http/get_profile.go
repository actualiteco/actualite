package profile_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"

	"go.uber.org/zap"
)

// GetUserProfile godoc
// @Summary Получить профиль пользователя
// @Description Получить профиль пользователя с постами по id пользователя
// @Tags Profile
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} ProfileResponse
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 403 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/{userID} [get]
// @Security BearerAuth
func (h *ProfileHTTPHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := httputils.GetUUIDPathValue(r, "userID")
	if err != nil {
		log.Info("get userID from path value", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	user, err := h.profileService.GetProfile(ctx, userID)
	if err != nil {
		log.Error("failed to get user profile", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(toResponse(user), http.StatusOK)
}
