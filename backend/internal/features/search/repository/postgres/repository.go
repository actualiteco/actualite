package search_repository

import (
	"actualiteco/actualite/internal/core/domain"
	postgres_pool "actualiteco/actualite/internal/core/repository/postgres"
	"fmt"
	"time"
	"uuid"
)

type SearchRepository struct {
	pool postgres_pool.Pool
}

func New(pool postgres_pool.Pool) *SearchRepository {
	return &SearchRepository{pool: pool}
}

type model struct {
	ID               uuid.UUID  `db:"id"`
	UserID           uuid.UUID  `db:"user_id"`
	Username         string     `db:"username"`
	DisplayName      *string    `db:"display_name"`
	AvatarURL        *string    `db:"avatar_url"`
	ParentID         *uuid.UUID `db:"parent_id"`
	Content          string     `db:"content"`
	ImageURL         *string    `db:"image_url"`
	LikesCount       uint       `db:"likes_count"`
	DeletedAt        *time.Time `db:"deleted_at"`
	CreatedAt        time.Time  `db:"created_at"`
	FromSubscription bool       `db:"from_subscription"`
	IsLiked          bool       `db:"is_liked"`
	Tags             []string   `db:"tags"`
}

func toPostDomain(m model) (domain.Post, error) {
	tags, err := domain.StringsToTags(m.Tags)
	if err != nil {
		return domain.Post{}, fmt.Errorf("parse tags: %w", err)
	}

	return domain.Post{
		ID:         m.ID,
		UserID:     m.UserID,
		ParentID:   m.ParentID,
		Content:    m.Content,
		ImageURL:   m.ImageURL,
		LikesCount: m.LikesCount,
		CreatedAt:  m.CreatedAt,
		IsLiked:    m.IsLiked,
		Hashtags:   tags,
	}, nil
}

func toPostWithAuthorDomain(m model) (domain.PostWithAuthor, error) {
	post, err := toPostDomain(m)
	if err != nil {
		return domain.PostWithAuthor{}, fmt.Errorf("post to domain: %w", err)
	}
	post.DeletedAt = m.DeletedAt

	author := domain.UserSummary{
		ID:          m.UserID,
		Username:    m.Username,
		DisplayName: m.DisplayName,
		AvatarURL:   m.AvatarURL,
	}

	return domain.NewPostWithAuthor(post, author), nil
}

func toPostWithAuthorDomains(models []model) ([]domain.PostWithAuthor, error) {
	posts := make([]domain.PostWithAuthor, len(models))
	for i, model := range models {
		post, err := toPostWithAuthorDomain(model)
		if err != nil {
			return nil, fmt.Errorf("to post with author domain: %w", err)
		}
		posts[i] = post
	}

	return posts, nil
}
