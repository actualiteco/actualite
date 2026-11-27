package search_repository

import (
	"context"
	"fmt"
	"strings"

	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
)

const userSearchLimit = 20

var escapeLike = strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_")

func (r *SearchRepository) SearchUsers(ctx context.Context, query string) ([]domain.UserSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const sql = `
		SELECT id::uuid, username, display_name, avatar_url
		FROM users
		WHERE username ILIKE $1 OR display_name ILIKE $1
		ORDER BY lower(username), id
		LIMIT $2
	`
	pattern := "%" + escapeLike.Replace(query) + "%"
	rows, err := postgres_pool.ExecutorFromContext(ctx, r.pool).Query(ctx, sql, pattern, userSearchLimit)
	if err != nil {
		return nil, fmt.Errorf("query matching users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.UserSummary, 0)
	for rows.Next() {
		var user domain.UserSummary
		if err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL); err != nil {
			return nil, fmt.Errorf("scan matching user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read matching users: %w", err)
	}

	return users, nil
}
