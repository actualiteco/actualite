package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *PostsRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (domain.PostWithAuthor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT p.id::uuid, p.user_id::uuid, u.username, u.display_name, u.avatar_url,
			p.parent_id::uuid, p.content, p.image_url, p.likes_count, p.created_at,
			p.deleted_at
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id 
		WHERE p.id = $1;
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, postID)

	var model PostModel
	err := row.Scan(&model.ID, &model.UserID, &model.Username, &model.DisplayName, &model.AvatarURL,
		&model.ParentID, &model.Content, &model.ImageURL, &model.LikesCount, &model.CreatedAt,
		&model.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PostWithAuthor{}, fmt.Errorf("no post with id=%s: %w", postID.String(), domain.ErrNotFound)
		}
		return domain.PostWithAuthor{}, fmt.Errorf("scan row: %w", err)
	}

	return toPostWithAuthorDomain(model), nil
}
