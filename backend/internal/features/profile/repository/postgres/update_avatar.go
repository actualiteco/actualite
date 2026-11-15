package users_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (r *UsersRepository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET avatar_url = $1
	WHERE id = $2
	`

	tag, err := r.ConnPool.Exec(ctxTimeout, query, avatarURL, id.String())
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no updated rows: %w", domain.ErrNotFound)
	}

	return nil
}
