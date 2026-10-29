package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tryingmyb3st/PolyTweet/internal/core/domain"
	hash_utils "github.com/tryingmyb3st/PolyTweet/internal/utils/hash"
)

func (s *AuthService) RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	user.ID = uuid.NewString()

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", domain.INVALID_REQUEST)
	}

	hash, err := hash_utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash generating: %w", domain.INTERNAL_ERROR)
	}

	user.Password = hash
	user.CreatedAt = time.Now()

	createdUser, err := s.authRepo.SaveNewUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("email is already booked: %w", domain.EMAIL_ALREADY_USED)
			}
		}

		return nil, fmt.Errorf("saving to database: %w", domain.INTERNAL_ERROR)
	}

	return createdUser, nil
}
