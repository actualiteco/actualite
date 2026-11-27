package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (r *PostsRepository) GetPostByID(ctx context.Context, userID, postID uuid.UUID) (domain.PostWithAuthor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
		SELECT p.id::uuid, p.user_id::uuid, u.username, u.display_name, u.avatar_url,
			p.parent_id::uuid, p.content, p.image_url, p.likes_count, p.created_at, p.deleted_at,
			EXISTS(
				SELECT 1 FROM likes l
				WHERE l.post_id = p.id
					AND l.user_id = $1
			) AS is_liked,
			parent.content AS parent_content,
			(SELECT COALESCE(array_agg(h.name ORDER BY h.name), '{}'::text[])
				FROM post_hashtags ph
				JOIN hashtags h ON h.id = ph.hashtag_id
				WHERE ph.post_id = p.id) AS tags
		FROM posts p
		LEFT JOIN users u      ON u.id      = p.user_id
		LEFT JOIN posts parent ON parent.id = p.parent_id
		WHERE p.id = $2;
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, userID, postID)

	var model PostModel
	err := row.Scan(&model.ID, &model.UserID, &model.Username, &model.DisplayName, &model.AvatarURL,
		&model.ParentID, &model.Content, &model.ImageURL, &model.LikesCount, &model.CreatedAt,
		&model.DeletedAt, &model.IsLiked, &model.ParentContent, &model.Tags)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PostWithAuthor{}, fmt.Errorf("no post with id=%s: %w", postID.String(), domain.ErrNotFound)
		}
		return domain.PostWithAuthor{}, fmt.Errorf("scan row: %w", err)
	}

	return toPostWithAuthorDomain(model)
}
