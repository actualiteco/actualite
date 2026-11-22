package auth_service

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
)

type AuthService struct {
	usersRepo    usersRepository
	sessionsRepo sessionsRepository
	tokenService tokenService
	mailer       mailer
	txManager    txManager

	frontendURL *url.URL
}

type usersRepository interface {
	SaveNewUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id, userID uuid.UUID) (domain.User, error)
	SavePasswordResetToken(ctx context.Context, userID uuid.UUID, email, tokenHash string, expirationTime time.Time) error
	GetUserByPasswordResetToken(ctx context.Context, tokenHash string) (domain.User, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	GetOwnProfile(ctx context.Context, id uuid.UUID) (domain.User, error)
	SaveEmailConfirmationToken(ctx context.Context, userID uuid.UUID, email, tokenHash string, expirationTime time.Time) error
	ConfirmEmail(ctx context.Context, userID uuid.UUID, tokenHash string) error
}

type sessionsRepository interface {
	Create(ctx context.Context, session domain.Session) (domain.Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	GetByRefreshHash(ctx context.Context, hash string) (domain.Session, error)
	RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, oldHash string, newHash string, expiresAt time.Time) error
	Revoke(ctx context.Context, userID, sessionID uuid.UUID) error
}

type tokenService interface {
	GenerateRefreshToken() (string, error)
	HashRefreshToken(token string) string
	GetRefreshTokenExpiration() time.Time

	GenerateJWT(userID, sessionID uuid.UUID, role string, isEmailVerified bool) (string, error)
	GetAccessTokenTTL() time.Duration

	GetPasswordResetTokenExpiration() time.Time
	GeneratePasswordResetToken() (string, error)
	HashPasswordResetToken(token string) string

	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool

	GetEmailConfirmationTokenExpiration() time.Time
	GenerateEmailConfirmationToken() (string, error)
	HashEmailConfirmationToken(token string) string
}

type mailer interface {
	SendLinkToResetPassword(ctx context.Context, to, link string) error
	GetSendEmailTimeout() time.Duration
	SendLinkToConfirmEmail(ctx context.Context, to, link string) error
}

type txManager interface {
	DoWrite(ctx context.Context, fn func(context.Context) error) error
}

func NewAuthService(usersRepo usersRepository, sessionsRepo sessionsRepository, tokenSvc tokenService,
	mailer mailer, tx txManager, frontendBaseURL string) (*AuthService, error) {
	u, err := url.Parse(frontendBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse frontend base url: %w", err)
	}

	return &AuthService{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		tokenService: tokenSvc,
		mailer:       mailer,
		txManager:    tx,
		frontendURL:  u,
	}, nil
}
