package auth_repository

import (
	"actualiteco/actualite/internal/core/domain"
	auth_models "actualiteco/actualite/internal/features/auth/repository"
	"context"
	"fmt"
)

func (r *AuthRepository) SaveNewUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO users(id, username, email, password, role, created_at)
	VALUES($1,$2,$3,$4,$5,$6)
	RETURNING id, username, email, password, role, created_at
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
	)

	var model auth_models.UserModel
	err := row.Scan(
		&model.ID,
		&model.Username,
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
		Username:  model.Username,
		Email:     model.Email,
		Password:  model.Password,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
	}, nil
}
