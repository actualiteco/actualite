package search_service

import (
	"context"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

type SearchRepository interface {
	SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error)
	CountPostsWithContent(ctx context.Context, query string) (uint64, error)
	SearchPostsByContent(ctx context.Context, query string, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
	CountPostsWithHashtags(ctx context.Context, tags []domain.Tag) (uint64, error)
	SearchPostsByHashtags(ctx context.Context, tags []domain.Tag, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
}

type txManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}

type SearchService struct {
	repo      SearchRepository
	txManager txManager
}

func New(repo SearchRepository, txManager txManager) *SearchService {
	return &SearchService{repo: repo, txManager: txManager}
}
