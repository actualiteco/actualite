package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"
)

func (s *AuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	tokenHash := s.tokenService.HashPasswordResetToken(token)

	err := s.txManager.DoWrite(ctx, func(ctx context.Context) error {
		user, err := s.usersRepo.GetUserByPasswordResetToken(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return fmt.Errorf("get password reset token: %v: %w", err, domain.ErrResetTokenInvalidOrExpired)
			}
			return fmt.Errorf("get password reset token: %w", err)
		}

		if time.Now().After(user.PasswordResetExpiresAt) {
			return domain.ErrResetTokenInvalidOrExpired
		}

		userID, err := uuid.Parse(user.ID)
		if err != nil {
			return fmt.Errorf("parse uuid=%s: %w", userID, err)
		}

		passwordHash, err := s.tokenService.HashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("hash new password: %w", err)
		}

		if err = s.usersRepo.ResetPassword(ctx, userID, passwordHash); err != nil {
			return fmt.Errorf("reset password: %w", err)
		}

		return nil
	})

	return err
}
