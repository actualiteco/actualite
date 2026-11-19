package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"context"
	"errors"
	"fmt"
	"uuid"

	"go.uber.org/zap"
)

func (s *AuthService) ProcessResetPasswordRequest(ctx context.Context, email string) error {
	user, err := s.usersRepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("get user: %w", err)
	}

	if !user.EmailVerified {
		return nil
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("parse userID=%s: %w", user.ID, err)
	}

	go func() {
		log := logger.FromContext(ctx)

		bgCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.mailer.GetSendEmailTimeout())
		defer cancel()

		if err := s.generateAndSendResetPasswordToken(bgCtx, userID, email); err != nil {
			log.Error(
				"generate and send reset password token",
				zap.String("userID", userID.String()),
				zap.Error(err),
			)
		}
	}()

	return nil
}

func (s *AuthService) generateAndSendResetPasswordToken(ctx context.Context, userID uuid.UUID, email string) error {
	token, err := s.tokenService.GeneratePasswordResetToken()
	if err != nil {
		return fmt.Errorf("generate reset password token: %w", err)
	}

	tokenHash := s.tokenService.HashPasswordResetToken(token)
	expirationTime := s.tokenService.GetPasswordResetTokenExpiration()
	if err := s.usersRepo.SavePasswordResetToken(ctx, userID, email, tokenHash, expirationTime); err != nil {
		return fmt.Errorf("save password reset token: %w", err)
	}

	if err = s.mailer.SendResetPasswordToken(ctx, email, token); err != nil {
		return fmt.Errorf("send reset password token email: %w", err)
	}

	return nil
}
