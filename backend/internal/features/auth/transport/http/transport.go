package auth_transport

import (
	"context"

	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"

	"github.com/tryingmyb3st/PolyTweet/internal/core/transport/server"
)

type AuthHTTPHandler struct {
	AuthService AuthService
}

type AuthService interface {
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)
}

func NewAuthHandler(authServ AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		AuthService: authServ,
	}
}

func (h *AuthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "POST",
			URL:     "/register",
			Handler: h.RegisterUser,
		},
	}
}
