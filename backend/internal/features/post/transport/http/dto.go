package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"time"
	"uuid"
)

type CreatePostDTO struct {
	Content  string     `json:"content" example:"Hello, world!" validate:"required,min=1,max=280"`
	ParentID *uuid.UUID `json:"parent_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"omitempty"`
	ImageURL *string    `json:"image_url" example:"https://example.com/image.jpg" validate:"omitempty,url"`
}

type PostDTOResponse struct {
	ID         uuid.UUID  `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	UserID     uuid.UUID  `json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	ParentID   *uuid.UUID `json:"parent_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content    string     `json:"content" example:"Hello, world!"`
	ImageURL   *string    `json:"image_url" example:"https://example.com/image.jpg"`
	LikesCount uint       `json:"likes_count" example:"10"`
	CreatedAt  time.Time  `json:"created_at" example:"2026-03-25T12:00:41.267Z"`
}

type PostWithAuthorDTOResponse struct {
	ID          uuid.UUID  `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	UserID      uuid.UUID  `json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Username    string     `json:"username" example:"testuser"`
	DisplayName *string    `json:"display_name" example:"test"`
	AvatarURL   *string    `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	ParentID    *uuid.UUID `json:"parent_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content     string     `json:"content" example:"Hello, world!"`
	ImageURL    *string    `json:"image_url" example:"https://example.com/image.jpg"`
	LikesCount  uint       `json:"likes_count" example:"10"`
	CreatedAt   time.Time  `json:"created_at" example:"2026-03-25T12:00:41.267Z"`
	DeletedAt   *time.Time `json:"deleted_at" example:"2026-03-25T12:00:41.267Z"`
}

type paginationDTO struct {
	Page       uint64 `json:"page" example:"1"`
	PageSize   uint64 `json:"page_size" example:"20"`
	Total      uint64 `json:"total" example:"20"`
	TotalPages uint64 `json:"total_pages" example:"100"`
}

type PostsWithPaginationDTOResponse struct {
	Posts      []PostWithAuthorDTOResponse `json:"posts"`
	Pagination paginationDTO               `json:"pagination"`
}

func postToResponse(p domain.Post) PostDTOResponse {
	return PostDTOResponse{
		ID:         p.ID,
		UserID:     p.UserID,
		ParentID:   p.ParentID,
		Content:    p.Content,
		ImageURL:   p.ImageURL,
		LikesCount: p.LikesCount,
		CreatedAt:  p.CreatedAt,
	}
}

func postWithAuthorToResponse(p domain.PostWithAuthor) PostWithAuthorDTOResponse {
	return PostWithAuthorDTOResponse{
		ID:          p.Post.ID,
		UserID:      p.Post.UserID,
		Username:    p.Author.Username,
		DisplayName: p.Author.DisplayName,
		AvatarURL:   p.Author.AvatarURL,
		ParentID:    p.Post.ParentID,
		Content:     p.Post.Content,
		ImageURL:    p.Post.ImageURL,
		LikesCount:  p.Post.LikesCount,
		CreatedAt:   p.Post.CreatedAt,
		DeletedAt:   p.Post.DeletedAt,
	}
}

func paginationToResponse(p domain.Pagination) paginationDTO {
	return paginationDTO{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      p.Total,
		TotalPages: p.TotalPages,
	}
}

func postsWithPaginationToResponse(posts []domain.PostWithAuthor, pagination domain.Pagination) PostsWithPaginationDTOResponse {
	postsResponse := make([]PostWithAuthorDTOResponse, len(posts))
	for i, post := range posts {
		postsResponse[i] = postWithAuthorToResponse(post)
	}

	paginationResponse := paginationToResponse(pagination)

	return PostsWithPaginationDTOResponse{
		Posts:      postsResponse,
		Pagination: paginationResponse,
	}
}
