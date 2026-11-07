package domain

import (
	"fmt"
	"time"
	"uuid"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	ID         uuid.UUID  `validate:"required"`
	UserID     uuid.UUID  `validate:"required"`
	ParentID   *uuid.UUID `validate:"omitempty"`
	Content    string     `validate:"required,min=1,max=280"`
	ImageURL   *string    `validate:"omitempty,url"`
	LikesCount uint
	DeletedAt  *time.Time
	CreatedAt  time.Time
}

type UserSummary struct {
	ID          uuid.UUID
	Username    string
	DisplayName *string
	AvatarURL   *string
}

type PostWithAuthor struct {
	Post   Post
	Author UserSummary
}

func NewPost(userID uuid.UUID, content string, imageURL *string, parentID *uuid.UUID) (Post, error) {
	post := Post{
		ID:        uuid.New(),
		UserID:    userID,
		ParentID:  parentID,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}

	if err := post.Validate(); err != nil {
		return Post{}, fmt.Errorf("failed to create new post: %w: %w", err, ErrInvalidRequest)
	}

	return post, nil
}

func NewPostWithAuthor(post Post, author UserSummary) PostWithAuthor {
	return PostWithAuthor{
		Post:   post,
		Author: author,
	}
}

func (p *Post) Validate() error {
	validator := validator.New()

	if err := validator.Struct(p); err != nil {
		return fmt.Errorf("validator struct: %w: %w", err, ErrInvalidRequest)
	}

	return nil
}
