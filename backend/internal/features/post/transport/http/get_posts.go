package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"
)

// GetPosts godoc
// @Summary Возвращает посты для ленты
// @Description Возвращаются посты подписок пользователя и посты за последний месяц (отсортированные по количеству лайков) поддержкой пагинации через параметры page и page_size
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"  minimum(1)  default(1)
// @Param page_size query int false "Размер страницы"  minimum(1)  maximum(30)  default(20)
// @Success 200 {object} PostsWithPaginationDTOResponse "Посты возвращены"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts [get]
// @Security BearerAuth
func (h *PostsHTTPHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	page, pageSize, err := httputils.GetPaginationParamsFromQuery(r)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	userID, err := middleware.GetUserUUIDFromCtx(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	posts, pagination, err := h.postsService.GetPosts(ctx, userID, page, pageSize)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(PostsWithPaginationToResponse(posts, pagination), http.StatusOK)
}
