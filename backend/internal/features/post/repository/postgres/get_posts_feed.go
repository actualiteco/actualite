package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *PostsRepository) GetPostsFeed(ctx context.Context, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		WITH user_follows AS (
			SELECT following_id AS author_id
			FROM follows
			WHERE follower_id = $1
		)
		SELECT
			p.id::uuid                       AS id,
			p.user_id::uuid                  AS user_id,
			u.username                       AS username,
			u.display_name                   AS display_name,
			u.avatar_url                     AS avatar_url,
			p.parent_id::uuid                AS parent_id,
			p.content                        AS content,
			p.image_url                      AS image_url,
			p.likes_count                    AS likes_count,
			p.created_at                     AS created_at,
			p.deleted_at                     AS deleted_at,
			(uf.author_id IS NOT NULL)       AS from_subscription,
			EXISTS(
				SELECT 1 FROM likes l
				WHERE l.post_id = p.id AND l.user_id = $1
			)                                AS is_liked,
			(SELECT COALESCE(array_agg(h.name ORDER BY h.name), '{}'::text[])
				FROM post_hashtags ph
				JOIN hashtags h ON h.id = ph.hashtag_id
				WHERE ph.post_id = p.id) 	 AS tags								
		FROM posts p
		JOIN users u              ON u.id         = p.user_id
		LEFT JOIN user_follows uf ON uf.author_id = p.user_id
		WHERE p.parent_id IS NULL
		AND p.deleted_at IS NULL
		AND (
				uf.author_id IS NOT NULL
			OR p.created_at >= NOW() - INTERVAL '1 month'
		)
		ORDER BY
			from_subscription DESC,
			CASE WHEN uf.author_id IS NOT NULL
				THEN EXTRACT(EPOCH FROM p.created_at)::bigint
				ELSE p.likes_count::bigint
			END DESC,
			p.id DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := conn.Query(ctxTimeout, query, userID, limit, offset)
	if err != nil {
		return []domain.PostWithAuthor{}, fmt.Errorf("get posts feed query: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[PostModel])
	if err != nil {
		return []domain.PostWithAuthor{}, fmt.Errorf("collect rows: %w", err)
	}

	return toPostWithAuthorDomains(models)
}
