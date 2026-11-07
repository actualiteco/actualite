package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/middleware"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

type ProfileHTTPHandler struct {
	profileService profileService
}

type profileService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, displayName, username, status *string) (domain.User, error)
}

func NewProfileHTTPHandler(profileService profileService) *ProfileHTTPHandler {
	return &ProfileHTTPHandler{
		profileService: profileService,
	}
}

func (h *ProfileHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:               http.MethodGet,
			URL:                  "/profile/me",
			Handler:              h.GetUserProfileByJWT,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/profile/{userID}",
			Handler:              h.GetUserProfile,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
		{
			Method:               http.MethodPatch,
			URL:                  "/profile",
			Handler:              h.Update,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware()},
		},
	}
}
