package tokenclaims

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	SessionID string `json:"sid"`
	Role      string `json:"role,omitempty"`

	jwt.RegisteredClaims
}

func (c *CustomClaims) IsActive() bool {
	return c.ExpiresAt.After(time.Now())
}
