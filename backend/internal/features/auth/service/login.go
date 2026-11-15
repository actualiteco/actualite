package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	hash_utils "actualiteco/actualite/internal/utils/hash"
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (s *AuthService) Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, int64, error) {
	user, err := s.usersRepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", 0, fmt.Errorf("no user: %w", domain.ErrUnauthorized)
		}

		return "", "", 0, fmt.Errorf("get from database: %w", domain.ErrInternal)
	}

	if !hash_utils.CheckPasswordHash(password, user.Password) {
		return "", "", 0, fmt.Errorf("wrong password: %w", domain.ErrUnauthorized)
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("parse userID=%s", user.ID)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshTokenHash := s.tokenService.HashRefreshToken(refreshToken)
	expirationTime := time.Now().Add(time.Minute * 10).UTC()

	session := domain.NewSession(userID, refreshTokenHash, expirationTime, userAgent, ipAddress)

	session, err = s.sessionsRepo.Create(ctx, session)
	if err != nil {
		return "", "", 0, fmt.Errorf("save session: %w", err)
	}

	accessToken, err := s.tokenService.GenerateJWT(userID.String(), session.ID, user.Role)
	if err != nil {
		return "", "", 0, fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, refreshToken, expirationTime.Unix(), nil
}
