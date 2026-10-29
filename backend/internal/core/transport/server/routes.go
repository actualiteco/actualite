package server

import (
	"net/http"

	"github.com/tryingmyb3st/PolyTweet/internal/core/middleware"
)

type Route struct {
	Method               string
	URL                  string
	Handler              http.HandlerFunc
	AdditionalMiddleware []middleware.Middleware
}
