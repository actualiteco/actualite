package middleware

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"actualiteco/actualite/internal/utils/jwt_utils"
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type key string

var (
	userID key = "userID"
	role   key = "role"
)

func GetUserIDFromCtx(ctx context.Context) string {
	userIDValue, ok := ctx.Value(userID).(string)
	if !ok {
		panic("no userID in context")
	}
	return userIDValue
}

func AuthMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			respWriter := response.NewResponseHandler(log, w)

			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer") {
				log.Debug("unathourized user", zap.String("Authorization", r.Header.Get("Authorization")))

				respWriter.MapError(domain.ErrUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt_utils.VerifyJWTtoken(token)
			if err != nil {
				log.Debug("invalid token", zap.Any("claims", claims))

				respWriter.MapError(domain.ErrInvalidRequest)
				return
			}
			log.Debug("user claims", zap.String("id", claims.UserId), zap.String("role", claims.Role))

			if claims.ExpiresAt.Before(time.Now()) {
				log.Debug("token expired", zap.Time("expires_at", claims.ExpiresAt.Time))

				respWriter.MapError(domain.ErrInvalidRequest)
				return
			}

			ctx := context.WithValue(r.Context(), userID, claims.UserId)
			ctx = context.WithValue(ctx, role, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
