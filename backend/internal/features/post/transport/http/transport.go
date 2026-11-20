package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

const (
	postIDPathValueKey    = "postID"
	userIDPathValueKey    = "userID"
	pageQueryValueKey     = "page"
	pageSizeQueryValueKey = "page_size"
)

type PostsHTTPHandler struct {
	postsService postsService
}

type postsService interface {
	Like(ctx context.Context, userID, postID uuid.UUID) error
	Unlike(ctx context.Context, userID, postID uuid.UUID) error
	Create(ctx context.Context, userID uuid.UUID, content string, parentID *uuid.UUID, imageURL *string) (domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error)
	GetPostsByAuthorID(ctx context.Context, userID, authorID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	GetReplies(ctx context.Context, userID, postID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	GetPosts(ctx context.Context, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error)
	Delete(ctx context.Context, userID, postID uuid.UUID) error
	UploadImage(ctx context.Context, data []byte) (string, error)
}

func NewPostsHandler(postsService postsService) *PostsHTTPHandler {
	return &PostsHTTPHandler{
		postsService: postsService,
	}
}

func (h *PostsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:       http.MethodPost,
			URL:          "/posts/{postID}/like",
			Handler:      h.Like,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/posts/{postID}/unlike",
			Handler:      h.Unlike,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/posts/create",
			Handler:      h.Create,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/posts/{postID}",
			Handler:      h.GetPostByID,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/profile/{userID}/posts",
			Handler:      h.GetPostsByUserID,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/posts/{postID}/replies",
			Handler:      h.GetReplies,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/posts",
			Handler:      h.GetPosts,
			AuthRequired: true,
		},
		{
			Method:       http.MethodDelete,
			URL:          "/posts/{postID}",
			Handler:      h.DeletePost,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/posts/me",
			Handler:      h.GetPostsByJWT,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/posts/image",
			Handler:      h.UploadImage,
			AuthRequired: true,
		},
	}
}
