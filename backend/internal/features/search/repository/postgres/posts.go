package search_repository

import (
	"context"
	"fmt"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"

	"github.com/jackc/pgx/v5"
)

const postSearchFilter = `
	WHERE p.deleted_at IS NULL
	  AND p.search_vector @@ plainto_tsquery('pg_catalog.simple', $1)
`

func (r *SearchRepository) CountPostsWithContent(ctx context.Context, query string) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const sql = "SELECT count(*) FROM posts p" + postSearchFilter
	var count uint64
	err := postgres_pool.ExecutorFromContext(ctx, r.pool).QueryRow(ctx, sql, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count matching posts: %w", err)
	}
	return count, nil
}

func (r *SearchRepository) SearchPostsByContent(ctx context.Context, query string, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const sql = `
		SELECT p.id::uuid, p.user_id::uuid, p.parent_id::uuid, p.content,
			p.image_url, p.likes_count, p.created_at, p.deleted_at,
			u.username, u.display_name, u.avatar_url,
			EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $2) AS is_liked,
			COALESCE(t.tags, '{}'::text[]) AS tags
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN LATERAL (
			SELECT array_agg(h.name ORDER BY h.name) AS tags
			FROM post_hashtags ph
			JOIN hashtags h ON h.id = ph.hashtag_id
			WHERE ph.post_id = p.id
		) t ON true
	` + postSearchFilter + `
		ORDER BY ts_rank(p.search_vector, plainto_tsquery('pg_catalog.simple', $1)) DESC,
				p.created_at DESC, p.id DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := postgres_pool.ExecutorFromContext(ctx, r.pool).Query(ctx, sql, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query matching posts: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return toPostWithAuthorDomains(models)
}

func (r *SearchRepository) CountPostsWithHashtags(ctx context.Context, tags []domain.Tag) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if len(tags) == 0 {
		return 0, nil
	}

	conn := postgres_pool.ExecutorFromContext(ctx, r.pool)

	query := `
		SELECT COUNT(DISTINCT ph.post_id)
		FROM post_hashtags ph
		JOIN hashtags h ON h.id = ph.hashtag_id
		WHERE h.name = ANY($1::text[]);
	`

	var count uint64
	if err := conn.QueryRow(ctx, query, tags).Scan(&count); err != nil {
		return 0, fmt.Errorf("count posts with hashtags: %w: %w", err, domain.ErrInternal)
	}

	return count, nil
}

func (r *SearchRepository) SearchPostsByHashtags(ctx context.Context, tags []domain.Tag,
	userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.pool)

	query := `
		SELECT p.id::uuid, p.user_id::uuid, u.username, u.display_name, u.avatar_url,
			p.parent_id::uuid, p.content, p.image_url, p.likes_count, p.created_at,
			EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS is_liked,
			p.deleted_at, COALESCE(t.tags, '{}'::text[]) AS tags
		FROM (
			SELECT DISTINCT ph.post_id
			FROM post_hashtags ph
			JOIN hashtags h ON h.id = ph.hashtag_id
			WHERE h.name = ANY($2::text[])
		) matched
		JOIN posts p ON p.id = matched.post_id
		JOIN users u ON u.id = p.user_id
		LEFT JOIN LATERAL (
			SELECT array_agg(h2.name ORDER BY h2.name) AS tags
			FROM post_hashtags ph2
			JOIN hashtags h2 ON h2.id = ph2.hashtag_id
			WHERE ph2.post_id = p.id
		) t ON true
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := conn.Query(ctx, query, userID, tags, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query matching posts: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[model])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return toPostWithAuthorDomains(models)
}
