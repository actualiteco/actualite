package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
)

const foreignKeyViolation = "23503"

func (r *UsersRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO follows (follower_id, following_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	tag, err := r.ConnPool.Exec(ctxTimeout, query, followerID.String(), followingID.String())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
			return fmt.Errorf("user with id=%s not found: %w", followingID.String(), domain.ErrNotFound)
		}
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("already following: %w", domain.ErrInvalidRequest)
	}

	return nil
}
