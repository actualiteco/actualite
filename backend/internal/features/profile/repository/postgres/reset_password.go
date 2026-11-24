package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"uuid"
)

func (r *UsersRepository) ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		UPDATE public.users
		SET password_reset_token = NULL, password_reset_expires_at = NULL, password = $1
		WHERE id = $2;
	`

	result, err := conn.Exec(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated: %w", domain.ErrNotFound)
	}

	return nil
}
