package auth_transport

import (
	"context"
	"net/http"
	"uuid"

	"actualiteco/actualite/internal/core/domain"

	"actualiteco/actualite/internal/core/transport/server"
)

type AuthHTTPHandler struct {
	AuthService AuthService
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.User, password string) (domain.User, error)
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, int64, error)
	Refresh(ctx context.Context, oldRefreshToken string) (string, string, int64, error)
	Logout(ctx context.Context, userID, sessionID uuid.UUID) error
}

func NewAuthHandler(authServ AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		AuthService: authServ,
	}
}

func (h *AuthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:       http.MethodPost,
			URL:          "/auth/register",
			Handler:      h.RegisterUser,
			AuthRequired: false,
		},
		{
			Method:       http.MethodPost,
			URL:          "/auth/login",
			Handler:      h.LoginUser,
			AuthRequired: false,
		},
		{
			Method:       http.MethodPost,
			URL:          "/auth/refresh",
			Handler:      h.Refresh,
			AuthRequired: false,
		},
		{
			Method:       http.MethodPost,
			URL:          "/auth/logout",
			Handler:      h.Logout,
			AuthRequired: true,
		},
	}
}
