package auth_repository

import (
	"context"
	"fmt"

	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"
	auth_models "github.com/tryingmyb3st/PolyTweet/internal/features/auth/repository"
)

func (r *AuthRepository) SaveNewUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO users(id, email, password, role, created_at)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, email, password, role, created_at
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		user.ID,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
	)

	var model auth_models.UserModel
	err := row.Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.Role,
		&model.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("scan returning user: %w", err)
	}

	return &domain.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
	}, nil
}
