package server

import (
	"net/http"

	"actualiteco/actualite/internal/core/middleware"
)

type Route struct {
	Method               string
	URL                  string
	Handler              http.HandlerFunc
	AdditionalMiddleware []middleware.Middleware
}
