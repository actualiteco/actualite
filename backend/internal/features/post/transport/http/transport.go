package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/middleware"
	tokenclaims "actualiteco/actualite/internal/core/token/service/claims"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

type PostsHTTPHandler struct {
	postsService postsService
	tokenService tokenService
}

type postsService interface {
	Like(ctx context.Context, userID, postID uuid.UUID) error
	Unlike(ctx context.Context, userID, postID uuid.UUID) error
	Create(ctx context.Context, userID uuid.UUID, content string, parentID *uuid.UUID, imageURL *string) (domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error)
	GetPostsByAuthorID(ctx context.Context, userID, authorID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	GetReplies(ctx context.Context, userID, postID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	GetPosts(ctx context.Context, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
}

type tokenService interface {
	VerifyJWT(tokenString string) (*tokenclaims.CustomClaims, error)
}

func NewPostsHandler(postsService postsService, tokenService tokenService) *PostsHTTPHandler {
	return &PostsHTTPHandler{
		postsService: postsService,
		tokenService: tokenService,
	}
}

func (h *PostsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:               http.MethodPost,
			URL:                  "/posts/{postID}/like",
			Handler:              h.Like,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodPost,
			URL:                  "/posts/{postID}/unlike",
			Handler:              h.Unlike,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodPost,
			URL:                  "/posts/create",
			Handler:              h.Create,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts/{postID}",
			Handler:              h.GetPostByID,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/profile/{userID}/posts",
			Handler:              h.GetPostsByUserID,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts/{postID}/replies",
			Handler:              h.GetReplies,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/posts",
			Handler:              h.GetPosts,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
	}
}
