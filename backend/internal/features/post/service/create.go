package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (s *PostsService) Create(ctx context.Context, userID uuid.UUID, content string, parentID *uuid.UUID, imageURL *string) (domain.Post, error) {
	post, err := domain.NewPost(userID, content, imageURL, parentID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("create post domain: %w", err)
	}

	createdPost, err := s.postsRepo.CreatePost(ctx, post)
	if err != nil {
		return domain.Post{}, fmt.Errorf("create post: %w: %w", err, domain.ErrInternal)
	}

	return createdPost, nil
}
