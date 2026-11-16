package auth_service

import (
	"context"
	"uuid"
)

func (s *AuthService) Logout(ctx context.Context, userID, sessionID uuid.UUID) error {
	return s.sessionsRepo.Revoke(ctx, userID, sessionID)
}
