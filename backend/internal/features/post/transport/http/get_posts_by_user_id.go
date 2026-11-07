package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"net/http"

	"go.uber.org/zap"
)

// GetPostsByUserID godoc
// @Summary Поиск постов пользователя
// @Description Ищет посты пользователя по ID с поддержкой пагинации через параметры page и page_size
// @Tags Posts
// @Accept json
// @Produce json
// @Param userID path string true "User ID (UUID)" Format(uuid) Example(3fa85f64-5717-4562-b3fc-2c963f66afa6)
// @Param page query int false "Номер страницы"  minimum(1)  default(1)
// @Param page_size query int false "Размер страницы"  minimum(1)  maximum(30)  default(20)
// @Success 200 {object} PostsWithPaginationDTOResponse "Посты найдены"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 404 {object} domain.CustomError "У пользователя нет постов"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/{userID}/posts [get]
// @Security BearerAuth
func (h *PostsHTTPHandler) GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
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

	userID, err := httputils.GetUUIDPathValue(r, "userID")
	if err != nil {
		respWriter.MapError(err)
		return
	}

	posts, pagination, err := h.PostsService.GetPostsByUserID(ctx, userID, page, pageSize)
	if err != nil {
		log.Error("get posts by user", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(postsWithPaginationToResponse(posts, pagination), http.StatusOK)
}
