package search_service

import (
	"context"

	"actualiteco/actualite/internal/core/domain"
)

func (s *SearchService) SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error) {
	query, err := trimQuery(query)
	if err != nil {
		return nil, err
	}
	return s.repo.SearchUsers(ctx, query)
}
