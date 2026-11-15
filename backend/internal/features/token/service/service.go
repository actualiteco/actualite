package tokenservice

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	config config
}

func New(cfg config) *TokenService {
	return &TokenService{
		config: cfg,
	}
}

type CustomClaims struct {
	SessionID string `json:"sid"`
	Role      string `json:"role,omitempty"`

	jwt.RegisteredClaims
}

func (t *TokenService) GenerateJWT(userID string, sessionID uuid.UUID, role string) (string, error) {
	expirationTime := time.Now().Add(t.config.AccessTTL).UTC()

	claims := CustomClaims{
		SessionID: sessionID.String(),
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "actualiteco",
			Audience:  []string{"actualiteco-api"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(t.config.Secret)
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}

	return tokenString, nil
}

func (t *TokenService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (t *TokenService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
