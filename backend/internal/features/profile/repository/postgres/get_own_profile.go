package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	auth_models "actualiteco/actualite/internal/features/profile/repository"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetOwnProfile(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, role, created_at, username, display_name, status, avatar_url, email_verified,
		(SELECT COUNT(*) FROM follows WHERE following_id = users.id) AS followers_count,
		(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
		email_confirmation_token, email_confirmation_expires_at
	FROM users
	WHERE id=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, id)

	var model auth_models.UserModel
	err := row.Scan(&model.ID, &model.Email, &model.Role, &model.CreatedAt, &model.Username, &model.DisplayName,
		&model.Status, &model.AvatarURL, &model.EmailVerified, &model.FollowersCount, &model.FollowingCount,
		&model.EmailConfirmationToken, &model.EmailConfirmationExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan returning user: %w", err)
	}

	return auth_models.ToDomain(model), nil
}
