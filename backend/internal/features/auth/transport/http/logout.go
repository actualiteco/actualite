package auth_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	"actualiteco/actualite/internal/core/transport/response"
	"context"
	"fmt"
	"net/http"
	"uuid"
)

// Logout godoc
// @Summary Инвалидация refresh-токена
// @Description Инвалидирует refresh-токен. Пока привязанный к нему access-токен не протух, по нему еще можно будет делать запросы (будет исправлено позже после интеграции Redis), но при это обновлять refresh-токен уже нельзя. Операция идемпотентна, можно несколько раз отозвать один и тот же refresh-токен без ошибки
// @Tags Auth
// @Accept json
// @Produce json
// @Success 204
// @Failure 401 {object} domain.CustomError "Access-токен протух или не валиден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /auth/logout [post]
// @Security BearerAuth
func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	sessionID, err := getSessionUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	if err := h.AuthService.Logout(ctx, userID, sessionID); err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.NoContent()
}

func getUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return uuid.Nil(), fmt.Errorf("no userID in context: %w", domain.ErrInternal)
	}

	return uuid.Parse(userID)
}

func getSessionUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	sessionID := middleware.GetSessionIDFromCtx(ctx)
	if sessionID == "" {
		return uuid.Nil(), fmt.Errorf("no sessionID in context: %w", domain.ErrInternal)
	}

	return uuid.Parse(sessionID)
}
