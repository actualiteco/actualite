package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	auth_models "actualiteco/actualite/internal/features/profile/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUser(ctx context.Context, email string) (domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, password, role, created_at, username, display_name, status, avatar_url, email_verified
	FROM users
	WHERE email=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, email)

	var model auth_models.UserModel
	err := row.Scan(&model.ID, &model.Email, &model.Password, &model.Role, &model.CreatedAt,
		&model.Username, &model.DisplayName, &model.Status, &model.AvatarURL, &model.EmailVerified)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("no user with email=%s: %w", email, domain.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan returning user: %w", err)
	}

	return auth_models.ToDomain(model), nil
}
