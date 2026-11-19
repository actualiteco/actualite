package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"time"
	"uuid"
)

func (r *UsersRepository) SavePasswordResetToken(ctx context.Context, userID uuid.UUID,
	email, tokenHash string, expirationTime time.Time) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		UPDATE public.users
		SET password_reset_token = $1, password_reset_expires_at = $2
		WHERE id = $3 AND email = $4;
	`

	tag, err := conn.Exec(ctxTimeout, query, tokenHash, expirationTime, userID, email)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated: %w", domain.ErrNotFound)
	}

	return nil
}
