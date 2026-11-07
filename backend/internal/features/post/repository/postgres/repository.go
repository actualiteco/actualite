package post_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"time"
	"uuid"
)

type PostsRepository struct {
	ConnPool postgres_pool.Pool
}

func New(connPool postgres_pool.Pool) *PostsRepository {
	return &PostsRepository{
		ConnPool: connPool,
	}
}

type PostModel struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Username    string     `db:"username"`
	DisplayName *string    `db:"display_name"`
	AvatarURL   *string    `db:"avatar_url"`
	ParentID    *uuid.UUID `db:"parent_id"`
	Content     string     `db:"content"`
	ImageURL    *string    `db:"image_url"`
	LikesCount  uint       `db:"likes_count"`
	DeletedAt   *time.Time `db:"deleted_at"`
	CreatedAt   time.Time  `db:"created_at"`
}

func toPostDomain(m PostModel) domain.Post {
	return domain.Post{
		ID:         m.ID,
		UserID:     m.UserID,
		ParentID:   m.ParentID,
		Content:    m.Content,
		ImageURL:   m.ImageURL,
		LikesCount: m.LikesCount,
		CreatedAt:  m.CreatedAt,
	}
}

func toPostWithAuthorDomain(m PostModel) domain.PostWithAuthor {
	post := toPostDomain(m)
	post.DeletedAt = m.DeletedAt

	author := domain.UserSummary{
		ID:          m.UserID,
		Username:    m.Username,
		DisplayName: m.DisplayName,
		AvatarURL:   m.AvatarURL,
	}

	return domain.NewPostWithAuthor(post, author)
}

func toPostWithAuthorDomains(models []PostModel) []domain.PostWithAuthor {
	p := make([]domain.PostWithAuthor, len(models))
	for i, model := range models {
		p[i] = toPostWithAuthorDomain(model)
	}
	return p
}
