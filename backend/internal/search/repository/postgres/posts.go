package search_repository

import (
	"context"
	"fmt"
	"uuid"

	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
)

const postSearchFilter = `
	WHERE p.deleted_at IS NULL
	  AND p.search_vector @@ plainto_tsquery('pg_catalog.simple', $1)
`

func (r *SearchRepository) CountPosts(ctx context.Context, query string) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const sql = "SELECT count(*) FROM posts p" + postSearchFilter
	var count uint64
	err := postgres_pool.ExecutorFromContext(ctx, r.pool).QueryRow(ctx, sql, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count matching posts: %w", err)
	}
	return count, nil
}

func (r *SearchRepository) SearchPosts(ctx context.Context, query string, userID uuid.UUID, limit, offset uint64) ([]domain.PostWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const sql = `
		SELECT p.id::uuid, p.user_id::uuid, p.parent_id::uuid, p.content,
		       p.image_url, p.likes_count, p.created_at, p.deleted_at,
		       u.username, u.display_name, u.avatar_url,
		       EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $2) AS is_liked
		FROM posts p
		JOIN users u ON u.id = p.user_id
	` + postSearchFilter + `
		ORDER BY ts_rank(p.search_vector, plainto_tsquery('pg_catalog.simple', $1)) DESC,
		         p.created_at DESC, p.id DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := postgres_pool.ExecutorFromContext(ctx, r.pool).Query(ctx, sql, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query matching posts: %w", err)
	}
	defer rows.Close()

	posts := make([]domain.PostWithAuthor, 0)
	for rows.Next() {
		var post domain.Post
		var author domain.UserSummary
		err := rows.Scan(
			&post.ID, &post.UserID, &post.ParentID, &post.Content,
			&post.ImageURL, &post.LikesCount, &post.CreatedAt, &post.DeletedAt,
			&author.Username, &author.DisplayName, &author.AvatarURL, &post.IsLiked,
		)
		if err != nil {
			return nil, fmt.Errorf("scan matching post: %w", err)
		}
		author.ID = post.UserID
		posts = append(posts, domain.NewPostWithAuthor(post, author))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read matching posts: %w", err)
	}

	return posts, nil
}
