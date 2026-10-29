package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"actualiteco/actualite/internal/core/domain"
	hash_utils "actualiteco/actualite/internal/utils/hash"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *AuthService) RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	user.ID = uuid.NewString()

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", domain.ErrInvalidRequest)
	}

	hash, err := hash_utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash generating: %w", domain.ErrInternal)
	}

	user.Password = hash
	user.CreatedAt = time.Now()

	createdUser, err := s.authRepo.SaveNewUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("email is already booked: %w", domain.ErrEmailAlreadyUsed)
			}
		}

		return nil, fmt.Errorf("saving to database: %w", domain.ErrInternal)
	}

	return createdUser, nil
}
