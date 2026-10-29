package auth_service

import (
	"context"

	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"
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
