package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"uuid"
)

func (r *PostsRepository) Delete(ctx context.Context, userID, postID uuid.UUID) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		UPDATE public.posts
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2;
	`

	tag, err := conn.Exec(ctxTimeout, query, postID, userID)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no posts deleted: %w", domain.ErrNotFound)
	}

	return nil
}
