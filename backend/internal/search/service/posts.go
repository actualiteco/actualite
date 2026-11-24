package search_service

import (
	"context"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

func (s *SearchService) SearchPosts(ctx context.Context, query string, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	query, err := trimQuery(query)
	if err != nil {
		return nil, domain.Pagination{}, err
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
