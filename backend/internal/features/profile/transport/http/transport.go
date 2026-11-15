package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/middleware"
	tokenclaims "actualiteco/actualite/internal/core/token/service/claims"
	"actualiteco/actualite/internal/core/transport/server"
	"context"
	"net/http"
	"uuid"
)

type ProfileHTTPHandler struct {
	profileService profileService
	tokenService   tokenService
}

type profileService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, displayName, username, status *string) (domain.User, error)
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsUsernameTaken(ctx context.Context, username string) (bool, error)
}

type tokenService interface {
	VerifyJWT(tokenString string) (*tokenclaims.CustomClaims, error)
}

func NewProfileHTTPHandler(profileService profileService, tokenService tokenService) *ProfileHTTPHandler {
	return &ProfileHTTPHandler{
		profileService: profileService,
		tokenService:   tokenService,
	}
}

func (h *ProfileHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:               http.MethodGet,
			URL:                  "/profile/me",
			Handler:              h.GetUserProfileByJWT,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodGet,
			URL:                  "/profile/{userID}",
			Handler:              h.GetUserProfile,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodPatch,
			URL:                  "/profile",
			Handler:              h.Update,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodPost,
			URL:                  "/profile/{userID}/follow",
			Handler:              h.Follow,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:               http.MethodPost,
			URL:                  "/profile/{userID}/unfollow",
			Handler:              h.Unfollow,
			AdditionalMiddleware: []middleware.Middleware{middleware.AuthMiddleware(h.tokenService)},
		},
		{
			Method:  http.MethodGet,
			URL:     "/profile/username-availability",
			Handler: h.CheckUsernameAvailability,
		},
	}
}
