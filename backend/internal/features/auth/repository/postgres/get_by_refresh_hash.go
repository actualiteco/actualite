package sessionsrepo

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *SessionsRepository) GetByRefreshHash(ctx context.Context, hash string) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, user_agent, COALESCE(ip_address::text, '')
        FROM sessions
        WHERE refresh_token_hash = $1;
	`

	row := r.ConnPool.QueryRow(ctx, query, hash)

	model, err := scanSession(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, fmt.Errorf("no session found: %v: %w", err, domain.ErrNotFound)
		}
		return domain.Session{}, err
	}

	return toDomain(model)
}
