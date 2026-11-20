package auth_service

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"context"
	"fmt"
	"uuid"

	"go.uber.org/zap"
)

func (s *AuthService) ProcessEmailConfirmationRequest(ctx context.Context, userID uuid.UUID) error {
	user, err := s.usersRepo.GetOwnProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user profile: %w", err)
	}

	if user.EmailVerified {
		return fmt.Errorf("user already verified: %w", domain.ErrAlreadyVerified)
	}

	go func() {
		log := logger.FromContext(ctx)

		bgCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.mailer.GetSendEmailTimeout())
		defer cancel()

		if err := s.generateAndSendEmailConfirmationToken(bgCtx, userID, user.Email); err != nil {
			log.Error(
				"generate and send email confirmation token",
				zap.String("userID", userID.String()),
				zap.String("email", user.Email),
				zap.Error(err),
			)
		} else {
			log.Debug("mail was successfully sent", zap.String("to", user.Email))
		}
	}()

	return nil
}

func (s *AuthService) generateAndSendEmailConfirmationToken(ctx context.Context, userID uuid.UUID, email string) error {
	token, err := s.tokenService.GenerateEmailConfirmationToken()
	if err != nil {
		return fmt.Errorf("generate email confirmation token: %w", err)
	}

	tokenHash := s.tokenService.HashEmailConfirmationToken(token)
	expirationTime := s.tokenService.GetEmailConfirmationTokenExpiration()

	link, err := createLinkWithToken(token)
	if err != nil {
		return fmt.Errorf("create link to confirm email: %w", err)
	}

	if err = s.mailer.SendLinkToConfirmEmail(ctx, email, link); err != nil {
		return fmt.Errorf("send confirmation email token: %w", err)
	}

	if err := s.usersRepo.SaveEmailConfirmationToken(ctx, userID, email, tokenHash, expirationTime); err != nil {
		return fmt.Errorf("save email confirmation token: %w", err)
	}

	return nil
}
