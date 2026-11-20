package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"time"
	"uuid"
)

func (r *UsersRepository) SaveEmailConfirmationToken(ctx context.Context, userID uuid.UUID, email, tokenHash string, expirationTime time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		UPDATE public.users
		SET email_confirmation_token = $1, email_confirmation_expires_at = $2
		WHERE id = $3 AND email = $4;
	`

	tag, err := conn.Exec(ctx, query, tokenHash, expirationTime, userID, email)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated: %w", domain.ErrNotFound)
	}

	return nil
}
