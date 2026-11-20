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

func (r *UsersRepository) GetUserByID(ctx context.Context, id, targetID uuid.UUID) (domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, role, username, display_name, status, avatar_url, email_verified,
		(SELECT COUNT(*) FROM follows WHERE following_id = users.id) AS followers_count,
		(SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count,
		EXISTS (
           SELECT 1
           FROM follows f
           WHERE f.follower_id = $1
             AND f.following_id = users.id
       ) AS is_following
	FROM users
	WHERE id=$2
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, id, targetID)

	var model auth_models.UserModel
	err := row.Scan(&model.ID, &model.Role, &model.Username, &model.DisplayName, &model.Status, &model.AvatarURL,
		&model.EmailVerified, &model.FollowersCount, &model.FollowingCount, &model.IsFollowing)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan returning user: %w", err)
	}

	return auth_models.ToDomain(model), nil
}
