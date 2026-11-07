package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	auth_models "actualiteco/actualite/internal/features/profile/repository"
	"context"
	"fmt"
	"uuid"
)

func (r *UsersRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, password, role, created_at, username, display_name, status, avatar_url, email_verified,
		(SELECT COUNT(*) FROM follows WHERE following_id = users.id) AS followers_count,
		(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count
	FROM users
	WHERE id=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, id)

	var model auth_models.UserModel
	err := row.Scan(&model.ID, &model.Email, &model.Password, &model.Role, &model.CreatedAt,
		&model.Username, &model.DisplayName, &model.Status, &model.AvatarURL, &model.EmailVerified,
		&model.FollowersCount, &model.FollowingCount)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan returning user: %w", err)
	}

	return auth_models.ToDomain(model), nil
}
