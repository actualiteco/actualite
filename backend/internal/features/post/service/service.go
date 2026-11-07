package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"uuid"
)

type PostsService struct {
	postsRepo PostsRepository
}

type PostsRepository interface {
	CreatePost(ctx context.Context, post domain.Post) (domain.Post, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error)
	GetPostsByUserID(ctx context.Context, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
	GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)

	GetCountOfPostsByUserID(ctx context.Context, userID uuid.UUID) (uint64, error)
	GetCountOfReplies(ctx context.Context, parentID uuid.UUID) (uint64, error)
}

func New(postsRepo PostsRepository) *PostsService {
	return &PostsService{
		postsRepo: postsRepo,
	}
}
