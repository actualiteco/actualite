package auth_service

import (
	"context"
	"time"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

type AuthService struct {
	usersRepo    usersRepository
	sessionsRepo sessionsRepository
	tokenService tokenService
}

type usersRepository interface {
	SaveNewUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, email string) (domain.User, error)
}

type sessionsRepository interface {
	Create(ctx context.Context, session domain.Session) (domain.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	GetByRefreshHash(ctx context.Context, hash string) (domain.Session, error)
	IsActive(ctx context.Context, id uuid.UUID) (bool, error)
	RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, oldHash string, newHash string, expiresAt time.Time) error
	Revoke(ctx context.Context, id uuid.UUID) error
}

type tokenService interface {
	GenerateRefreshToken() (string, error)
	HashRefreshToken(token string) string
	GenerateJWT(userID, sessionID uuid.UUID, role string) (string, error)
	GetRefreshTokenExpiration() time.Time
	GetAccessTokenTTL() time.Duration
}

func NewAuthService(usersRepo usersRepository, sessionsRepo sessionsRepository, tokenSvc tokenService) *AuthService {
	return &AuthService{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		tokenService: tokenSvc,
	}
}
