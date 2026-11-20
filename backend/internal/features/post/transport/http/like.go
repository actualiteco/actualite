package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"

	"go.uber.org/zap"
)

type LikeResponse struct {
	Message string `json:"message"`
}

// Like godoc
// @Summary Поставить лайк
// @Description Поставить лайк на пост авторизованному пользователю
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path string true "Post ID"
// @Success 200 {object} LikeResponse "Лайк успешно поставлен"
// @Failure 400 {object} domain.CustomError "Невалидный post ID или состояние лайка"
// @Failure 401 {object} domain.CustomError "Пользователь не авторизован"
// @Failure 404 {object} domain.CustomError "Пост не найден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/{postID}/like [post]
// @Security BearerAuth
func (h *PostsHTTPHandler) Like(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)
	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}
	postID, err := httputils.GetUUIDPathValue(r, postIDPathValueKey)
	if err != nil {
		respWriter.MapError(err)
		return
	}
	if err := h.postsService.Like(ctx, userID, postID); err != nil {
		log.Error("like post", zap.Error(err))
		respWriter.MapError(err)
		return
	}
	respWriter.JSONResponse(LikeResponse{Message: "success"}, http.StatusOK)
}

// Unlike godoc
// @Summary Убрать лайк
// @Description Убрать лайк с поста авторизованному пользователю
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path string true "Post ID"
// @Success 200 {object} LikeResponse "Лайк успешно снят"
// @Failure 400 {object} domain.CustomError "Невалидный post ID или состояние лайка"
// @Failure 401 {object} domain.CustomError "Пользователь не авторизован"
// @Failure 404 {object} domain.CustomError "Пост не найден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/{postID}/unlike [post]
// @Security BearerAuth
func (h *PostsHTTPHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)
	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}
	postID, err := httputils.GetUUIDPathValue(r, postIDPathValueKey)
	if err != nil {
		respWriter.MapError(err)
		return
	}
	if err := h.postsService.Unlike(ctx, userID, postID); err != nil {
		log.Error("unlike post", zap.Error(err))
		respWriter.MapError(err)
		return
	}
	respWriter.JSONResponse(LikeResponse{Message: "success"}, http.StatusOK)
}
