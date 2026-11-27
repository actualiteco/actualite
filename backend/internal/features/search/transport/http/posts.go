package search_transport

import (
	"fmt"
	"net/http"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	post_transport "actualiteco/actualite/internal/features/post/transport/http"
	httputils "actualiteco/actualite/internal/utils/http"
)

type SearchPostsByQueryRequest struct {
	Query string `json:"query" validate:"required,min=1"`
}

// SearchPostsByQuery godoc
// @Summary Поиск постов по query (контенту)
// @Description Возвращает страницу результатов полнотекстового поиска по содержимому постов и ответов в порядке релевантности.
// @Tags Search
// @Produce json
// @Param request body SearchPostsByQueryRequest true "Текст для поиска"
// @Param page query int false "Номер страницы" minimum(1) default(1)
// @Param page_size query int false "Размер страницы" minimum(1) maximum(30) default(20)
// @Success 200 {object} post_transport.PostsWithPaginationDTOResponse "Найденные посты и пагинация"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Требуется авторизация"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /search/posts [post]
// @Security BearerAuth
func (h *SearchHTTPHandler) SearchPostsByQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	respWriter := response.NewResponseHandler(logger.FromContext(ctx), w)

	var req SearchPostsByQueryRequest
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		respWriter.MapError(err)
		return
	}

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

	userID, err := uuid.Parse(middleware.GetUserIDFromCtx(ctx))
	if err != nil {
		respWriter.MapError(fmt.Errorf("invalid user ID in token: %w", domain.ErrUnauthorized))
		return
	}

	posts, pagination, err := h.service.SearchPosts(ctx, req.Query, userID, page, pageSize)
	if err != nil {
		respWriter.MapError(err)
		return
	}
	respWriter.JSONResponse(post_transport.PostsWithPaginationToResponse(posts, pagination), http.StatusOK)
}
