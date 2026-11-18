package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/jackc/pgx/v5"
)

func (s *ProfileService) GetUserProfile(ctx context.Context, id, targetID uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id, targetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("no user with id=%s: %w", targetID.String(), domain.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
