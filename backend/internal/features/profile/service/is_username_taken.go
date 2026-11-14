package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"strings"
)

func (s *ProfileService) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	username = strings.ToLower(strings.TrimSpace(username))

	if err := domain.ValidateUsername(username); err != nil {
		return false, err
	}

	return s.repo.IsUsernameTaken(ctx, username)
}
