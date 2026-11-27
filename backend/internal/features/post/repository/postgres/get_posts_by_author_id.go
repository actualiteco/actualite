package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *PostsRepository) GetPostsByAuthorID(ctx context.Context, userID, authorID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctxTimeout, r.ConnPool)

	query := `
		SELECT p.id::uuid, p.user_id::uuid, u.username, u.display_name, u.avatar_url,
			p.parent_id::uuid, p.content, p.image_url, p.likes_count, p.created_at,
			p.deleted_at, 
			EXISTS(
				SELECT 1 FROM likes l
				WHERE l.post_id = p.id
				AND l.user_id = $1
			) AS is_liked,
			(SELECT COALESCE(array_agg(h.name ORDER BY h.name), '{}'::text[])
				FROM post_hashtags ph
				JOIN hashtags h ON h.id = ph.hashtag_id
				WHERE ph.post_id = p.id) AS tags
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.user_id = $2
		AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC, p.id DESC
		OFFSET $3
		LIMIT $4;
	`

	rows, err := conn.Query(ctxTimeout, query, userID, authorID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("query posts by author id: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[PostModel])
	if err != nil {
		return nil, fmt.Errorf("collect rows to model: %w", err)
	}

	if len(models) == 0 {
		return []domain.PostWithAuthor{}, nil
	}

	return toPostWithAuthorDomains(models)
}
