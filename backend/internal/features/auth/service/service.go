package auth_service

import (
	"context"

	"actualiteco/actualite/internal/core/domain"
)

type AuthService struct {
	authRepo AuthRepository
}

type AuthRepository interface {
	SaveNewUser(ctx context.Context, user domain.User) (*domain.User, error)
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		authRepo: repo,
	}
}
