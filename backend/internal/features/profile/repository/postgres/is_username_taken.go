package users_repository

import (
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"fmt"
)

func (r *UsersRepository) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	conn := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE LOWER(username) = $1
		);
	`

	row := conn.QueryRow(ctxTimeout, query, username)

	var usernameInUse bool
	if err := row.Scan(&usernameInUse); err != nil {
		return false, fmt.Errorf("scan row: %w", err)
	}

	return usernameInUse, nil
}
