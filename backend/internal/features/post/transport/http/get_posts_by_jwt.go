package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

// GetPostsByJWT godoc
// @Summary Поиск постов пользователя по его JWT
// @Description Ищет посты пользователя по его Access Token с поддержкой пагинации через параметры page и page_size
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"  minimum(1)  default(1)
// @Param page_size query int false "Размер страницы"  minimum(1)  maximum(30)  default(20)
// @Success 200 {object} PostsWithPaginationDTOResponse "Все существующие посты возвращены (если у пользователя нет постов - будет пустой массив)"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 404 {object} domain.CustomError "Пользователь не найден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/me [get]
// @Security BearerAuth
func (h *PostsHTTPHandler) GetPostsByJWT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	page, pageSize, err := getPaginationParamsFromQuery(r)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	posts, pagination, err := h.postsService.GetPostsByAuthorID(ctx, userID, userID, page, pageSize)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(postsWithPaginationToResponse(posts, pagination), http.StatusOK)
}
