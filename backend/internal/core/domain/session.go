package domain

import (
	"time"
	"uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string

	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time

	UserAgent string
	IPAddress string
}

func NewSession(userID uuid.UUID, refreshTokenHash string, expiresAt time.Time, userAgent, ipAddress string) Session {
	return Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        expiresAt,
		UserAgent:        userAgent,
		IPAddress:        ipAddress,
	}
}

func (s *Session) IsActive() bool {
	return s.RevokedAt == nil && time.Now().Before(s.ExpiresAt)
}
