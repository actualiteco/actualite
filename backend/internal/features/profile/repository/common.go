package auth_models

import (
	"actualiteco/actualite/internal/core/domain"
	"database/sql"
	"time"
)

type UserModel struct {
	ID                         string         `db:"id"`
	Email                      string         `db:"email"`
	Password                   string         `db:"password"`
	Role                       string         `db:"role"`
	CreatedAt                  time.Time      `db:"created_at"`
	Username                   string         `db:"username"`
	DisplayName                sql.NullString `db:"display_name"`
	AvatarURL                  sql.NullString `db:"avatar_url"`
	Status                     sql.NullString `db:"status"`
	EmailVerified              bool           `db:"email_verified"`
	FollowersCount             int            `db:"followers_count"`
	FollowingCount             int            `db:"following_count"`
	IsFollowing                bool           `db:"is_following"`
	PasswordResetToken         sql.NullString `db:"password_reset_token"`
	PasswordResetExpiresAt     sql.NullTime   `db:"password_reset_expires_at"`
	EmailConfirmationToken     sql.NullString `db:"email_confirmation_token"`
	EmailConfirmationExpiresAt sql.NullTime   `db:"email_confirmation_expires_at"`
}

func ToDomain(model UserModel) domain.User {
	user := domain.User{
		ID:                         model.ID,
		Email:                      model.Email,
		Password:                   model.Password,
		Role:                       model.Role,
		CreatedAt:                  model.CreatedAt,
		Username:                   model.Username,
		DisplayName:                model.DisplayName.String,
		AvatarURL:                  model.AvatarURL.String,
		Status:                     model.Status.String,
		EmailVerified:              model.EmailVerified,
		FollowersCount:             model.FollowersCount,
		FollowingCount:             model.FollowingCount,
		IsFollowing:                model.IsFollowing,
		PasswordResetToken:         model.PasswordResetToken.String,
		PasswordResetExpiresAt:     model.PasswordResetExpiresAt.Time,
		EmailConfirmationToken:     model.EmailConfirmationToken.String,
		EmailConfirmationExpiresAt: model.EmailConfirmationExpiresAt.Time,
	}

	return user
}
