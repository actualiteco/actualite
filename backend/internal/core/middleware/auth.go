package middleware

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	tokenclaims "actualiteco/actualite/internal/core/token/service/claims"
	"actualiteco/actualite/internal/core/transport/response"
	"context"
	"fmt"
	"net/http"
	"strings"
	"uuid"

	"go.uber.org/zap"
)

type key string

var (
	userIDKey        key = "userID"
	roleKey          key = "role"
	sessionIDKey     key = "sessionID"
	emailVerifiedKey key = "emailVerified"
)

type tokenService interface {
	VerifyJWT(tokenString string) (*tokenclaims.CustomClaims, error)
}

func GetUserUUIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	userIDValue, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil(), fmt.Errorf("no userID in context: %w", domain.ErrInternal)
	}

	return userIDValue, nil
}

func GetSessionUUIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	sessionIDValue, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil(), fmt.Errorf("no sessionID in context: %w", domain.ErrInternal)
	}

	return sessionIDValue, nil
}

func getEmailVerifiedFromCtx(ctx context.Context) bool {
	emailVerifiedValue, ok := ctx.Value(emailVerifiedKey).(bool)
	if !ok {
		panic("no emailVerified in context")
	}
	return emailVerifiedValue
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
				respWriter.MapError(err)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				log.Debug("invalid userID", zap.String("userID", claims.Subject))
				respWriter.MapError(fmt.Errorf("invalid userID: %w", domain.ErrInternal))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			ctx = context.WithValue(ctx, sessionIDKey, claims.SessionID)
			ctx = context.WithValue(ctx, emailVerifiedKey, claims.EmailVerified)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func EmailVerificationMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			respWriter := response.NewResponseHandler(log, w)

			if !getEmailVerifiedFromCtx(ctx) {
				err := fmt.Errorf("email is not verified: %w", domain.ErrEmailVerificationRequired)
				respWriter.MapError(err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
