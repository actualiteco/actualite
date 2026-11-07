package post_repository

import (
	"context"
	"fmt"
	"uuid"
)

func (r *PostsRepository) GetCountOfReplies(ctx context.Context, parentID uuid.UUID) (uint64, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT COUNT(*)
		FROM posts
		WHERE parent_id = $1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, parentID)

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get count: %w", err)
	}

	return count, nil
}
