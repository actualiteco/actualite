package auth_service

import (
	"context"
	"fmt"
	"time"

	"actualiteco/actualite/internal/core/domain"
	hash_utils "actualiteco/actualite/internal/utils/hash"

	"uuid"
)

func (s *AuthService) RegisterUser(ctx context.Context, user domain.User, password string) (domain.User, error) {
	user.ID = uuid.New().String()

	hash, err := hash_utils.HashPassword(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash generating: %w", domain.ErrInternal)
	}

	user.Password = hash
	user.CreatedAt = time.Now()

	createdUser, err := s.authRepo.SaveNewUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("saving to database: %w", err)
	}

	return createdUser, nil
}
