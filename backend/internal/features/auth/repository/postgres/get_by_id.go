package sessionsrepo

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"uuid"
)

func (r *SessionsRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, user_agent, COALESCE(ip_address::text, '')
        FROM sessions
        WHERE id = $1;
	`

	row := r.ConnPool.QueryRow(ctx, query, id)

	model, err := scanSession(row)
	if err != nil {
		return domain.Session{}, err
	}

	return toDomain(model)
}
