package search_service

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

func (s *SearchService) SearchPosts(ctx context.Context, query string, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	if query = strings.TrimSpace(query); query == "" {
		return nil, domain.Pagination{}, fmt.Errorf("query is required: %w", domain.ErrInvalidRequest)
	}

	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	var posts []domain.PostWithAuthor
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		count, err := s.repo.CountPosts(ctx, query)
		if err != nil {
			return err
		}
		pagination.ParseCount(count)
		if pagination.Page > pagination.TotalPages {
			return nil
		}

		posts, err = s.repo.SearchPosts(ctx, query, userID, limit, offset)
		return err
	})
	return posts, pagination, err
}
