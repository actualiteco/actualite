package search_transport

import (
	"uuid"
)

const (
	postIDPathValueKey = "postID"
	userIDPathValueKey = "userID"
	hashtagsQueryKey   = "hashtags"
)

type UserSearchDTOResponse struct {
	ID          uuid.UUID `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
}

type SearchPostsByQueryRequest struct {
	Query string `json:"query" validate:"required,min=1"`
}
