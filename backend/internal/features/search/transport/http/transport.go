package search_transport

import (
	"context"
	"net/http"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/transport/server"
)

type searchService interface {
	SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error)
	SearchPosts(ctx context.Context, query string, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
}

type SearchHTTPHandler struct {
	service searchService
}

func NewSearchHTTPHandler(service searchService) *SearchHTTPHandler {
	return &SearchHTTPHandler{service: service}
}

func (h *SearchHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:       http.MethodGet,
			URL:          "/search/users",
			Handler:      h.SearchUsers,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/search/posts",
			Handler:      h.SearchPostsByQuery,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/search/posts",
			AuthRequired: true,
		},
	}
}
