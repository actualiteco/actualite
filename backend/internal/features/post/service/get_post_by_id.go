package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"uuid"
)

func (s *PostsService) GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error) {
	return s.postsRepo.GetPostByID(ctx, postID)
}
