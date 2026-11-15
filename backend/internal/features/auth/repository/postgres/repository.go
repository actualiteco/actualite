package sessionsrepo

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"fmt"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type SessionsRepository struct {
	ConnPool postgres_pool.Pool
}

func New(connPool postgres_pool.Pool) *SessionsRepository {
	return &SessionsRepository{
		ConnPool: connPool,
	}
}

type model struct {
	ID               uuid.UUID  `db:"id"`
	UserID           string     `db:"user_id"`
	RefreshTokenHash string     `db:"refresh_token_hash"`
	CreatedAt        time.Time  `db:"created_at"`
	ExpiresAt        time.Time  `db:"expires_at"`
	RevokedAt        *time.Time `db:"revoked_at"`
	UserAgent        string     `db:"user_agent"`
	IPAddress        string     `db:"ip_address"`
}

func toDomain(m model) (domain.Session, error) {
	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("parse userID=%s: %w", m.UserID, err)
	}
	return domain.Session{
		ID:               m.ID,
		UserID:           userID,
		RefreshTokenHash: m.RefreshTokenHash,
		CreatedAt:        m.CreatedAt,
		ExpiresAt:        m.ExpiresAt,
		RevokedAt:        m.RevokedAt,
		UserAgent:        m.UserAgent,
		IPAddress:        m.IPAddress,
	}, nil
}

func scanSession(row pgx.Row) (model, error) {
	var m model
	err := row.Scan(&m.ID, &m.UserID, &m.RefreshTokenHash, &m.CreatedAt, &m.ExpiresAt, &m.RevokedAt, &m.UserAgent, &m.IPAddress)
	if err != nil {
		return model{}, fmt.Errorf("scan session: %w", err)
	}

	return m, nil
}
