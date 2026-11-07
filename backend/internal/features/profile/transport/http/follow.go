package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"

	"go.uber.org/zap"
)

// Follow godoc
// @Summary Подписаться на пользователя
// @Description Подписаться авторизованным пользователем на другого пользователя по userID
// @Tags Profile
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param Authorization header string true "Bearer <jwt токен>"
// @Success 200 {object} FollowResponse "Успешный запрос"
// @Failure 400 {object} domain.CustomError "Неверный запрос (в т.ч. подписка на себя или повторная подписка)"
// @Failure 401 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 404 {object} domain.CustomError "Пользователь не найден"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/{userID}/follow [post]
func (h *ProfileHTTPHandler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	followerID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Error("get userID from context", zap.Error(err))
		respWriter.ErrorResponse(domain.ErrInternal, http.StatusInternalServerError)
		return
	}

	followingID, err := httputils.GetUUIDPathValue(r, "userID")
	if err != nil {
		log.Info("get userID from path value", zap.Error(err))
		respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	if err := h.service.Follow(ctx, followerID, followingID); err != nil {
		log.Error("failed to follow user", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(FollowResponse{Message: "success"}, http.StatusOK)
}

// Unfollow godoc
// @Summary Отписаться от пользователя
// @Description Отписаться авторизованным пользователем от другого пользователя по userID
// @Tags Profile
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param Authorization header string true "Bearer <jwt токен>"
// @Success 200 {object} FollowResponse "Успешный запрос"
// @Failure 400 {object} domain.CustomError "Неверный запрос (в т.ч. отписка от себя или отсутствие подписки)"
// @Failure 401 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/{userID}/unfollow [post]
func (h *ProfileHTTPHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	followerID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Error("get userID from context", zap.Error(err))
		respWriter.ErrorResponse(domain.ErrInternal, http.StatusInternalServerError)
		return
	}

	followingID, err := httputils.GetUUIDPathValue(r, "userID")
	if err != nil {
		log.Info("get userID from path value", zap.Error(err))
		respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	if err := h.service.Unfollow(ctx, followerID, followingID); err != nil {
		log.Error("failed to unfollow user", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(FollowResponse{Message: "success"}, http.StatusOK)
}
