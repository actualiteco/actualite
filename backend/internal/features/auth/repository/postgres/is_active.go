package sessionsrepo

import (
	"context"
	"fmt"
	"uuid"
)

func (r *SessionsRepository) IsActive(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT EXISTS (
            SELECT 1
            FROM sessions
            WHERE id = $1
              AND revoked_at IS NULL
              AND expires_at > NOW()
        )
	`

	var isActive bool
	err := r.ConnPool.QueryRow(ctx, query, id).Scan(&isActive)
	if err != nil {
		return false, fmt.Errorf("check session activity: %w", err)
	}

	return isActive, nil
}
