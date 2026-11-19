package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (s *ProfileService) GetOwnProfile(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetOwnProfile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("no user with id=%s: %w", id.String(), domain.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
