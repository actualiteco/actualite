package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"uuid"
)

func (s *PostsService) GetPostByID(ctx context.Context, userID, postID uuid.UUID) (domain.PostWithAuthor, error) {
	return s.postsRepo.GetPostByID(ctx, userID, postID)
}
