package middleware

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	tokenclaims "actualiteco/actualite/internal/core/token/service/claims"
	"actualiteco/actualite/internal/core/transport/response"
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type key string

var (
	userIDKey    key = "userID"
	roleKey      key = "role"
	sessionIDKey key = "sessionID"
)

type tokenService interface {
	VerifyJWT(tokenString string) (*tokenclaims.CustomClaims, error)
}

func GetUserIDFromCtx(ctx context.Context) string {
	userIDValue, ok := ctx.Value(userIDKey).(string)
	if !ok {
		panic("no userID in context")
	}
	return userIDValue
}

func AuthMiddleware(tokenService tokenService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			respWriter := response.NewResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				respWriter.MapError(domain.ErrUnauthorized)
				return
			}

			claims, err := tokenService.VerifyJWT(parts[1])
			if err != nil {
				log.Debug("invalid token", zap.Any("claims", claims), zap.Error(err))
				respWriter.MapError(domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			ctx = context.WithValue(ctx, sessionIDKey, claims.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
