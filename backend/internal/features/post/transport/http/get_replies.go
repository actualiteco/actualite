package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"

	"go.uber.org/zap"
)

// GetReplies godoc
// @Summary Поиск комментариев на пост
// @Description Ищет комментарии на пост по его ID с поддержкой пагинации через параметры page и page_size
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path string true "Post ID (UUID)" Format(uuid) Example(a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d)
// @Param page query int false "Номер страницы"  minimum(1)  default(1)
// @Param page_size query int false "Размер страницы"  minimum(1)  maximum(30)  default(20)
// @Success 200 {object} PostsWithPaginationDTOResponse "Посты найдены"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 404 {object} domain.CustomError "У поста нет комментариев"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/{postID}/replies [get]
// @Security BearerAuth
func (h *PostsHTTPHandler) GetReplies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	page, err := httputils.GetUintQueryParam(r, "page")
	if err != nil {
		respWriter.MapError(err)
		return
	}

	pageSize, err := httputils.GetUintQueryParam(r, "page_size")
	if err != nil {
		respWriter.MapError(err)
		return
	}

	postID, err := httputils.GetUUIDPathValue(r, "postID")
	if err != nil {
		respWriter.MapError(err)
		return
	}

	log.Debug("getting replies for post", zap.String("postID", postID.String()))

	posts, pagination, err := h.PostsService.GetReplies(ctx, postID, page, pageSize)
	if err != nil {
		log.Error("get replies", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(postsWithPaginationToResponse(posts, pagination), http.StatusOK)
}
