package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

// TODO: wrap in transaction
func (s *PostsService) GetReplies(ctx context.Context, postID uuid.UUID, page, pageSize *uint64) ([]domain.PostWithAuthor, domain.Pagination, error) {
	pagination := domain.NewPagination(page, pageSize)
	offset, limit, err := pagination.GetOffsetLimit()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get offset and limit: %w", err)
	}

	count, err := s.postsRepo.GetCountOfReplies(ctx, postID)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get count of posts: %w", err)
	}
	if count == 0 {
		return nil, domain.Pagination{}, fmt.Errorf("no replies for post with id=%s: %w", postID.String(), domain.ErrNotFound)
	}

	pagination.ParseCount(count)

	posts, err := s.postsRepo.GetReplies(ctx, postID, limit, offset)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("get replies: %w", err)
	}

	return posts, pagination, nil
}
