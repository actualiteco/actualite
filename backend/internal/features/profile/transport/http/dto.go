package profile_transport

import "actualiteco/actualite/internal/core/domain"

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" example:"mark"`
	Username    *string `json:"username" example:"username"`
	Status      *string `json:"status" example:"something smart"`
}

type ProfileResponse struct {
	ID             string `json:"id" example:"8794e589-0ddb-43ce-9f92-16faafcf4ee4"`
	Email          string `json:"email" example:"test.user@email.com"`
	Role           string `json:"role" example:"user"`
	CreatedAt      string `json:"created_at" example:"2006-01-02T15:04:05Z07:00"`
	DisplayName    string `json:"display_name" example:"mark"`
	Username       string `json:"username" example:"username"`
	AvatarURL      string `json:"avatar_url" example:"http://localhost:8080/avatar_url"`
	Status         string `json:"status" example:"something smart"`
	EmailVerified  bool   `json:"email_verified" example:"false"`
	FollowersCount int    `json:"followers_count" example:"10"`
	FollowingCount int    `json:"following_count" example:"5"`
}

type FollowResponse struct {
	Message string `json:"message" example:"success"`
}

func toResponse(user domain.User) ProfileResponse {
	return ProfileResponse{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		AvatarURL:      user.AvatarURL,
		Status:         user.Status,
		EmailVerified:  user.EmailVerified,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
	}
}
