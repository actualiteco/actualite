package profile_transport

import "actualiteco/actualite/internal/core/domain"

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" example:"mark"`
	Username    *string `json:"username" example:"username"`
	Status      *string `json:"status" example:"something smart"`
}

type UploadAvatarResponse struct {
	AvatarURL string `json:"avatar_url" example:"http://localhost:8333/3,01637037d6"`
}

type MyProfileResponse struct {
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

type ProfileResponse struct {
	ID             string `json:"id" example:"8794e589-0ddb-43ce-9f92-16faafcf4ee4"`
	Role           string `json:"role" example:"user"`
	DisplayName    string `json:"display_name" example:"mark"`
	Username       string `json:"username" example:"username"`
	AvatarURL      string `json:"avatar_url" example:"http://localhost:8080/avatar_url"`
	Status         string `json:"status" example:"something smart"`
	FollowersCount int    `json:"followers_count" example:"10"`
	FollowingCount int    `json:"following_count" example:"5"`
	IsFollowing    bool   `json:"is_following" example:"true"`
}

type FollowResponse struct {
	Message string `json:"message" example:"success"`
}

type CheckUsernameAvailabilityResponse struct {
	Available bool `json:"available" example:"true"`
}

func toMyProfileResponse(user domain.User) MyProfileResponse {
	return MyProfileResponse{
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

func toProfileResponse(user domain.User) ProfileResponse {
	return ProfileResponse{
		ID:             user.ID,
		Role:           user.Role,
		Username:       user.Username,
		DisplayName:    user.DisplayName,
		AvatarURL:      user.AvatarURL,
		Status:         user.Status,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
		IsFollowing:    user.IsFollowing,
	}
}
