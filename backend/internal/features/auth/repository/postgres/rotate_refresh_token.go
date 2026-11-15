package sessionsrepo

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"time"
	"uuid"
)

func (r *SessionsRepository) RotateRefreshToken(ctx context.Context, sessionID uuid.UUID,
	oldHash string, newHash string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		UPDATE sessions
        SET refresh_token_hash = $1,
            expires_at = $2
        WHERE id = $3
          AND refresh_token_hash = $4
          AND revoked_at IS NULL
          AND expires_at > NOW()
	`

	result, err := r.ConnPool.Exec(ctx, query, newHash, expiresAt, sessionID, oldHash)
	if err != nil {
		return fmt.Errorf("rotate refresh token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("refresh token is invalid or expired: %w", domain.ErrUnauthorized)
	}

	return nil
}
