package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (s *PostsService) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := s.postsRepo.GetPostByID(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("get post by id: %w", err)
	}

	if post.Author.ID != userID {
		return fmt.Errorf("user is not the author of the post: %w", domain.ErrConflict)
	}

	if err := s.postsRepo.Delete(ctx, userID, postID); err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}
