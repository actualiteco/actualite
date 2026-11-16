package sessionsrepo

import (
	"context"
	"fmt"
	"uuid"
)

func (r *SessionsRepository) Revoke(ctx context.Context, userID, sessionID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		UPDATE sessions
        SET revoked_at = NOW()
        WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL;
	`

	_, err := r.ConnPool.Exec(ctx, query, sessionID, userID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}
