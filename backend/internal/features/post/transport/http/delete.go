package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"
)

// DeletePost godoc
// @Summary Удаление поста
// @Description Удаляет пост, если запрос на удаление исходит от автора поста. Метод идемпотентен
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path string true "Post ID (UUID)" Format(uuid) Example(3fa85f64-5717-4562-b3fc-2c963f66afa6)
// @Success 204 "Пост успешно удален"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 404 {object} domain.CustomError "Пост не найден"
// @Failure 409 {object} domain.CustomError "Пост принадлежит не текущему пользователю"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/{postID} [delete]
// @Security BearerAuth
func (h *PostsHTTPHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.postsService.Delete(ctx, userID, postID)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.NoContent()
}
