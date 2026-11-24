package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
)

func (s *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (string, string, int64, error) {
	oldRefreshTokenHash := s.tokenService.HashRefreshToken(oldRefreshToken)

	session, err := s.sessionsRepo.GetByRefreshHash(ctx, oldRefreshTokenHash)
	if err != nil {
		return "", "", 0, fmt.Errorf("get by refresh hash: %w", err)
	}

	if !session.IsActive() {
		return "", "", 0, fmt.Errorf("session is inactive, should login again: %w", domain.ErrSessionRevokedOrExpired)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, fmt.Errorf("generate new refresh token: %w", err)
	}

	newRefreshTokenHash := s.tokenService.HashRefreshToken(newRefreshToken)
	expirationTime := s.tokenService.GetRefreshTokenExpiration()

	user, err := s.usersRepo.GetUserByID(ctx, session.UserID, session.UserID)
	if err != nil {
		return "", "", 0, fmt.Errorf("get user by id: %w", err)
	}

	accessToken, err := s.tokenService.GenerateJWT(session.UserID, session.ID, user.Role, user.EmailVerified)
	if err != nil {
		return "", "", 0, fmt.Errorf("generate access token: %w", err)
	}

	if err := s.sessionsRepo.RotateRefreshToken(ctx, session.ID, oldRefreshTokenHash, newRefreshTokenHash, expirationTime); err != nil {
		return "", "", 0, fmt.Errorf("rotate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, int64(s.tokenService.GetAccessTokenTTL().Seconds()), nil
}
