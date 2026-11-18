package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

type ProfileHTTPHandler struct {
	profileService profileService
}

type profileService interface {
	GetOwnProfile(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserProfile(ctx context.Context, id, targetID uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, displayName, username, status *string) (domain.User, error)
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsUsernameTaken(ctx context.Context, username string) (bool, error)
	UploadAvatar(ctx context.Context, id uuid.UUID, data []byte) (string, error)
}

func NewProfileHTTPHandler(profileService profileService) *ProfileHTTPHandler {
	return &ProfileHTTPHandler{
		profileService: profileService,
	}
}

func (h *ProfileHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:       http.MethodGet,
			URL:          "/profile/me",
			Handler:      h.GetUserProfileByJWT,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/profile/{userID}",
			Handler:      h.GetUserProfile,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPatch,
			URL:          "/profile",
			Handler:      h.Update,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/profile/avatar",
			Handler:      h.UploadAvatar,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/profile/{userID}/follow",
			Handler:      h.Follow,
			AuthRequired: true,
		},
		{
			Method:       http.MethodPost,
			URL:          "/profile/{userID}/unfollow",
			Handler:      h.Unfollow,
			AuthRequired: true,
		},
		{
			Method:       http.MethodGet,
			URL:          "/profile/username-availability",
			Handler:      h.CheckUsernameAvailability,
			AuthRequired: false,
		},
	}
}
