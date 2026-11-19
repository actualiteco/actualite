package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	auth_models "actualiteco/actualite/internal/features/profile/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUserByPasswordResetToken(ctx context.Context, tokenHash string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		SELECT id, password_reset_token, password_reset_expires_at
		FROM public.users
		WHERE password_reset_token = $1;
	`

	row := conn.QueryRow(ctx, query, tokenHash)

	var model auth_models.UserModel
	if err := row.Scan(&model.ID, &model.PasswordResetToken, &model.PasswordResetExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("no row found: %w", domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan row: %w", err)
	}

	return auth_models.ToDomain(model), nil
}
