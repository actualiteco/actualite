package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"actualiteco/actualite/internal/utils/jwt_utils"

	"go.uber.org/zap"
)

var skipLogPrefixes = []string{"/swagger/", "/swagger"}

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func shouldSkipLogging(r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return true
	}

	path := r.URL.Path

	for _, p := range skipLogPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func LogMiddleware(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldSkipLogging(r) {
				next.ServeHTTP(w, r)
				return
			}

			log := log.With(zap.String("URL", r.URL.String()))

			ctx := context.WithValue(r.Context(), "log", log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PanicMiddleware(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := response.NewResponseHandler(log, w)

			defer func() {
				if err := recover(); err != nil {
					log.Error("panic happened", zap.Any("err", err))
					handler.PanicResponse(err, "panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func TraceMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldSkipLogging(r) {
				next.ServeHTTP(w, r)
				return
			}

			log := r.Context().Value("log").(*logger.Logger)

			before := time.Now()
			log.Info(">>> incoming request", zap.String("METHOD", r.Method))

			next.ServeHTTP(w, r)

			after := time.Now()

			log.Info("<<< request handled", zap.Duration("latency", time.Duration(after.Sub(before))))
		})
	}
}

func AuthMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*logger.Logger)
			respWriter := response.NewResponseHandler(log, w)

			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer") {
				log.Debug("unathourized user", zap.String("Authorization", r.Header.Get("Authorization")))

				respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt_utils.VerifyJWTtoken(token)
			if err != nil {
				log.Debug("invalid token", zap.Any("claims", claims))

				respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
				return
			}
			log.Debug("user claims", zap.String("id", claims.UserId), zap.String("role", claims.Role))

			if claims.ExpiresAt.Before(time.Now()) {
				log.Debug("token expired", zap.Time("expiresAt", claims.ExpiresAt.Time))

				respWriter.ErrorResponse(domain.ErrInvalidRequest, http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", claims.UserId)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CORSMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
