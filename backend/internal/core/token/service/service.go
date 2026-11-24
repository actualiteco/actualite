package tokenservice

import (
	"actualiteco/actualite/internal/core/domain"
	tokenclaims "actualiteco/actualite/internal/core/token/service/claims"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	config config
}

func New(cfg config) *TokenService {
	return &TokenService{
		config: cfg,
	}
}

func (s *TokenService) GenerateJWT(userID, sessionID uuid.UUID, role string, emailVerified bool) (string, error) {
	expirationTime := time.Now().Add(s.config.AccessTTL).UTC()

	claims := tokenclaims.CustomClaims{
		SessionID:     sessionID.String(),
		Role:          role,
		EmailVerified: emailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "actualiteco",
			Audience:  []string{"actualiteco-api"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.config.Secret)
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}

	return tokenString, nil
}

func (s *TokenService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (s *TokenService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func (s *TokenService) VerifyJWT(tokenString string) (*tokenclaims.CustomClaims, error) {
	claims := tokenclaims.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return s.config.Secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", domain.ErrAccessTokenExpired)
		}
		return nil, fmt.Errorf("error parsing claims: %v: %w", err, domain.ErrUnauthorized)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid: %v: %w", err, domain.ErrUnauthorized)
	}

	tokenClaims := token.Claims.(*tokenclaims.CustomClaims)

	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, fmt.Errorf("invalid userID=%s in token: %w", tokenClaims.Subject, domain.ErrUnauthorized)
	}

	if _, err := uuid.Parse(claims.SessionID); err != nil {
		return nil, fmt.Errorf("invalid sessionID=%s in token: %w", tokenClaims.SessionID, domain.ErrUnauthorized)
	}

	return tokenClaims, nil
}

func (s *TokenService) GetRefreshTokenExpiration() time.Time {
	return time.Now().Add(s.config.RefreshTTL).UTC()
}

func (s *TokenService) GetAccessTokenTTL() time.Duration {
	return s.config.AccessTTL
}

func (s *TokenService) GetPasswordResetTokenExpiration() time.Time {
	return time.Now().Add(s.config.PasswordResetTTL).UTC()
}

func (s *TokenService) GeneratePasswordResetToken() (string, error) {
	return s.GenerateRefreshToken()
}

func (s *TokenService) HashPasswordResetToken(token string) string {
	return s.HashRefreshToken(token)
}

func (s *TokenService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *TokenService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *TokenService) GetEmailConfirmationTokenExpiration() time.Time {
	return time.Now().Add(s.config.EmailConfirmationTTL).UTC()
}

func (s *TokenService) GenerateEmailConfirmationToken() (string, error) {
	return s.GenerateRefreshToken()
}

func (s *TokenService) HashEmailConfirmationToken(token string) string {
	return s.HashRefreshToken(token)
}
