package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
	"uuid"
)

func (r *UsersRepository) ConfirmEmail(ctx context.Context, userID uuid.UUID, tokenHash string) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		UPDATE public.users
		SET email_confirmation_token = NULL, email_confirmation_expires_at = NULL, email_verified = true
		WHERE id = $1 AND email_confirmation_token = $2;	
	`

	tag, err := conn.Exec(ctx, query, userID, tokenHash)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated: %w", domain.ErrNotFound)
	}

	return nil
}
