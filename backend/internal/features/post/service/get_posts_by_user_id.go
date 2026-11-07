package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

// TODO: wrap in transaction
func (s *PostsService) GetPostsByUserID(ctx context.Context, userID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get offset and limit: %w", err)
	}

	count, err := s.postsRepo.GetCountOfPostsByUserID(ctx, userID)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get count of posts: %w", err)
	}
	if count == 0 {
		return nil, domain.Pagination{}, fmt.Errorf("no posts from user with id=%s: %w", userID.String(), domain.ErrNotFound)
	}

	pagination.ParseCount(count)

	posts, err := s.postsRepo.GetPostsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get posts by userID: %w", err)
	}

	return posts, pagination, nil
}
