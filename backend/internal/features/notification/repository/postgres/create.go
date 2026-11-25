package notification_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5/pgconn"
)

const foreignKeyViolation = "23503"

func (r *NotificationsRepository) Create(ctx context.Context, n domain.Notification) error {
	ctx, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO notifications (id, recipient_id, actor_id, type, post_id, comment_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	executor := postgres_pool.ExecutorFromContext(ctx, r.ConnPool)
	_, err := executor.Exec(ctx, query,
		n.ID.String(),
		n.RecipientID.String(),
		n.ActorID.String(),
		string(n.Type),
		n.PostID.String(),
		uuidPtrToString(n.PostID),
		uuidPtrToString(n.CommentID),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolation {
			return fmt.Errorf("notification reference not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
