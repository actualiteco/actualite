package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/middleware"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

type PostsHTTPHandler struct {
	PostsService PostsService
}

type PostsService interface {
	Create(ctx context.Context, userID uuid.UUID, content string, parentID *uuid.UUID, imageURL *string) (domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error)
	GetPostsByUserID(ctx context.Context, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	GetReplies(ctx context.Context, postID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
}

func NewPostsHandler(postsService PostsService) *PostsHTTPHandler {
	return &PostsHTTPHandler{
		PostsService: postsService,
	}
}

func (h *PostsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:               http.MethodPost,
			URL:                  "/posts/create",
			Handler:              h.Create,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts/{postID}",
			Handler:              h.GetPostByID,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/profile/{userID}/posts",
			Handler:              h.GetPostsByUserID,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts/{postID}/replies",
			Handler:              h.GetReplies,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts",
			Handler:              h.GetPosts,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
	}
}
