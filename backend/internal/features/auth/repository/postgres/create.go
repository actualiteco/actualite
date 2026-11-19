package sessionsrepo

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
)

func (r *SessionsRepository) Create(ctx context.Context, session domain.Session) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO sessions (
            id,
            user_id,
            refresh_token_hash,
            expires_at,
            user_agent,
            ip_address
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, user_agent, ip_address::text;
	`

	row := r.ConnPool.QueryRow(ctx, query, session.ID, session.UserID, session.RefreshTokenHash,
		session.ExpiresAt, session.UserAgent, session.IPAddress)

	model, err := scanSession(row)
	if err != nil {
		return domain.Session{}, err
	}

	return toDomain(model)
}
