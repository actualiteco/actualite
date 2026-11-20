package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"time"
	"uuid"
)

func (s *AuthService) ConfirmEmail(ctx context.Context, userID uuid.UUID, token string) error {
	err := s.txManager.DoWrite(ctx, func(ctx context.Context) error {
		user, err := s.usersRepo.GetOwnProfile(ctx, userID)
		if err != nil {
			return fmt.Errorf("get user profile: %w", err)
		}

		tokenHash := s.tokenService.HashEmailConfirmationToken(token)

		if err := validateEmailConfirmationToken(user, tokenHash); err != nil {
			return fmt.Errorf("validate email confirmation token: %w", err)
		}

		if err := s.usersRepo.ConfirmEmail(ctx, userID, tokenHash); err != nil {
			return fmt.Errorf("confirm email: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func validateEmailConfirmationToken(user domain.User, tokenHash string) error {
	if time.Now().UTC().After(user.EmailConfirmationExpiresAt) {
		return fmt.Errorf("token expired: %w", domain.ErrTokenInvalidOrExpired)
	}

	if tokenHash != user.EmailConfirmationToken {
		return fmt.Errorf("invalid token: %w", domain.ErrTokenInvalidOrExpired)
	}

	return nil
}
