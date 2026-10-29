package auth_transport

import (
	"context"
	"net/http"

	"actualiteco/actualite/internal/core/domain"

	"actualiteco/actualite/internal/core/transport/server"
)

type AuthHTTPHandler struct {
	AuthService AuthService
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)
	LoginUser(ctx context.Context, email, password string) (*string, *domain.User, error)
}

func NewAuthHandler(authServ AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		AuthService: authServ,
	}
}

func (h *AuthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			URL:     "/register",
			Handler: h.RegisterUser,
		},
		{
			Method:  http.MethodPost,
			URL:     "/login",
			Handler: h.LoginUser,
		},
	}
}
