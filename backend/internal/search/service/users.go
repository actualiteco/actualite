package search_service

import (
	"context"
	"fmt"
	"strings"

	"actualiteco/actualite/internal/core/domain"
)

func (s *SearchService) SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error) {
	if query = strings.TrimSpace(query); query == "" {
		return nil, fmt.Errorf("query is required: %w", domain.ErrInvalidRequest)
	}
	return s.repo.SearchUsers(ctx, query)
}
