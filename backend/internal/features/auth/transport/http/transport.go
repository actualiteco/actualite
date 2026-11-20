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
	ProcessResetPasswordRequest(ctx context.Context, mail string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
	ProcessEmailConfirmationRequest(ctx context.Context, userID uuid.UUID) error
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
		{
			Method:  http.MethodPost,
			URL:     "/auth/password-reset/request",
			Handler: h.ProcessResetPasswordRequest,
		},
		{
			Method:  http.MethodPost,
			URL:     "/auth/password-reset/confirm",
			Handler: h.ConfirmPasswordReset,
		},
		{
			Method:       http.MethodPost,
			URL:          "/auth/email-confirmation/request",
			Handler:      h.ProcessEmailConfirmationRequest,
			AuthRequired: true,
		},
	}
}
