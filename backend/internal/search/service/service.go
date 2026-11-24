package search_service

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

type SearchRepository interface {
	SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error)
	CountPosts(ctx context.Context, query string) (uint64, error)
	SearchPosts(ctx context.Context, query string, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error)
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

func trimQuery(query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required: %w", domain.ErrInvalidRequest)
	}
	return query, nil
}
