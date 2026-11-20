package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const foreignKeyViolation = "23503"

func (r *PostsRepository) Like(ctx context.Context, userID, postID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	executor := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)
	var id string
	err := executor.QueryRow(ctx, `SELECT id FROM posts WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`, postID.String()).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("post not found: %w", domain.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("lock post: %w", err)
	}

	tag, err := executor.Exec(ctx, `INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT (user_id, post_id) DO NOTHING`, userID.String(), postID.String())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
			return fmt.Errorf("like reference not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("like post: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("invalid like state: %w", domain.ErrInvalidRequest)
	}

	if _, err := executor.Exec(ctx, `UPDATE posts SET likes_count = likes_count + 1 WHERE id = $1`, postID.String()); err != nil {
		return fmt.Errorf("update likes count: %w", err)
	}

	return nil
}
