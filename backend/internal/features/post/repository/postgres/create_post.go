package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *PostsRepository) CreatePost(ctx context.Context, post domain.Post, tags []domain.Tag) (domain.Post, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	tagNames := make([]string, len(tags))
	for i, t := range tags {
		tagNames[i] = string(t)
	}

	const query = `
		WITH input_tags AS (
			SELECT name, ord
			FROM unnest($6::text[]) WITH ORDINALITY AS t(name, ord)
		),
		upserted AS (
			INSERT INTO hashtags (id, name)
			SELECT gen_random_uuid()::text, name FROM input_tags
			ON CONFLICT (name) DO NOTHING
			RETURNING id, name
		),
		all_tags AS (
			SELECT u.id, i.name, i.ord
			FROM upserted u
			JOIN input_tags i ON i.name = u.name
			UNION ALL
			SELECT h.id, h.name, i.ord
			FROM hashtags h
			JOIN input_tags i ON h.name = i.name
			WHERE NOT EXISTS (SELECT 1 FROM upserted u WHERE u.name = h.name)
		),
		inserted_post AS (
			INSERT INTO posts (id, user_id, parent_id, content, image_url)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, user_id, parent_id, content, image_url,
					likes_count, deleted_at, created_at
		),
		inserted_links AS (
			INSERT INTO post_hashtags (post_id, hashtag_id)
			SELECT ip.id, at.id
			FROM inserted_post ip
			CROSS JOIN all_tags at
			RETURNING 1
		)
		SELECT ip.id::uuid,
			ip.user_id::uuid,
			ip.parent_id::uuid,
			ip.content,
			ip.image_url,
			ip.likes_count,
			ip.deleted_at,
			ip.created_at,
			(SELECT COALESCE(array_agg(name ORDER BY ord), '{}'::text[]) FROM all_tags) AS tags
		FROM inserted_post ip;
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query,
		post.ID.String(),
		post.UserID.String(),
		nullableUUID(post.ParentID),
		post.Content,
		post.ImageURL,
		tagNames,
	)

	var model PostModel
	err := row.Scan(
		&model.ID, &model.UserID, &model.ParentID, &model.Content, &model.ImageURL,
		&model.LikesCount, &model.DeletedAt, &model.CreatedAt, &model.Tags,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
			return domain.Post{}, fmt.Errorf(
				"failed to find parent_id=%s: %w: %w",
				post.ParentID, err, domain.ErrNotFound,
			)
		}
		return domain.Post{}, fmt.Errorf("scan returning post: %w", err)
	}

	return toPostDomain(model)
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}
	return id.String()
}
