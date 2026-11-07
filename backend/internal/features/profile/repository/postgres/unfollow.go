package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (r *UsersRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM follows
	WHERE follower_id = $1 AND following_id = $2
	`

	tag, err := r.ConnPool.Exec(ctxTimeout, query, followerID.String(), followingID.String())
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("not following: %w", domain.ErrInvalidRequest)
	}

	return nil
}
