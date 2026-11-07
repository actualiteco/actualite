package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *PostsRepository) GetReplies(ctx context.Context, parentID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT p.id::uuid, p.user_id::uuid, u.username, u.display_name, u.avatar_url,
			p.parent_id::uuid, p.content, p.image_url, p.likes_count, p.created_at,
			p.deleted_at
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id 
		WHERE p.parent_id = $1
		ORDER BY p.created_at DESC
		OFFSET $2
		LIMIT $3;
	`

	rows, err := r.ConnPool.Query(ctxTimeout, query, parentID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("query posts by user id: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[PostModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows to model: %w", err)
	}

	return toPostWithAuthorDomains(models), nil
}
