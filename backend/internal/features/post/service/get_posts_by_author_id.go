package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (s *PostsService) GetPostsByAuthorID(ctx context.Context, userID, authorID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get offset and limit: %w", err)
	}

	posts := make([]domain.PostWithAuthor, 0)
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		count, userExists, err := s.postsRepo.GetCountOfPostsByUserID(ctx, authorID)
		if err != nil {
			return fmt.Errorf("get count of posts: %w", err)
		}
		if !userExists {
			return fmt.Errorf("no user with id=%s: %w", authorID.String(), domain.ErrNotFound)
		}

		pagination.ParseCount(count)

		if count == 0 {
			return nil
		}

		posts, err = s.postsRepo.GetPostsByAuthorID(ctx, userID, authorID, limit, offset)
		if err != nil {
			return fmt.Errorf("get posts by authorID: %w", err)
		}

		return nil
	})
	if err != nil {
		return []domain.PostWithAuthor{}, domain.Pagination{}, err
	}

	return posts, pagination, nil
}
