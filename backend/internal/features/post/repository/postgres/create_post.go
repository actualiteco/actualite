package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *PostsRepository) CreatePost(ctx context.Context, post domain.Post) (domain.Post, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := sq.
		Insert("posts").
		Columns("id", "user_id", "parent_id", "content", "image_url").
		Values(post.ID, post.UserID, post.ParentID, post.Content, post.ImageURL).
		Suffix("RETURNING id::uuid, user_id::uuid, parent_id::uuid, content, image_url, likes_count, deleted_at, created_at").
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return domain.Post{}, fmt.Errorf("create sql query: %w: %w", err, domain.ErrInternal)
	}

	row := r.ConnPool.QueryRow(ctxTimeout, sqlQuery, args...)

	var model PostModel
	err = row.Scan(&model.ID, &model.UserID, &model.ParentID, &model.Content, &model.ImageURL,
		&model.LikesCount, &model.DeletedAt, &model.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.Post{}, fmt.Errorf("failed to find parent_id=%s: %w: %w", post.ParentID, err, domain.ErrNotFound)
		}
		return domain.Post{}, fmt.Errorf("scan returning post: %w", err)
	}

	return toPostDomain(model), nil
}
