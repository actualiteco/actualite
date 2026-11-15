package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	"actualiteco/actualite/internal/core/transport/response"
	"context"
	"fmt"
	"net/http"
	"uuid"

	"go.uber.org/zap"
)

// GetUserProfileByJWT godoc
// @Summary Получить профиль пользователя по JWT
// @Description Получить профиль пользователя с постами по JWT
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} ProfileResponse
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 403 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/me [get]
// @Security BearerAuth
func (h *ProfileHTTPHandler) GetUserProfileByJWT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Error("get userID from context", zap.Error(err))
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

func getUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return uuid.Nil(), fmt.Errorf("no userID in context: %w", domain.ErrInternal)
	}

	return uuid.Parse(userID)
}
