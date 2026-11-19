package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"fmt"
	"uuid"
)

func (s *ProfileService) Update(ctx context.Context, id uuid.UUID, displayName, username, status *string) (domain.User, error) {
	if displayName == nil && username == nil && status == nil {
		return domain.User{}, fmt.Errorf("no fields to update: %w", domain.ErrInvalidRequest)
	}

	var updatedUser domain.User
	err := s.txManager.DoWrite(ctx, func(ctx context.Context) error {
		user, err := s.repo.GetOwnProfile(ctx, id)
		if err != nil {
			return fmt.Errorf("get user: %w", err)
		}

		if displayName != nil {
			if err := user.UpdateDisplayName(*displayName); err != nil {
				return fmt.Errorf("update display_name: %w", err)
			}
		}

		if username != nil {
			if err := user.UpdateUsername(*username); err != nil {
				return fmt.Errorf("update username: %w", err)
			}
		}

		if status != nil {
			if err := user.UpdateStatus(*status); err != nil {
				return fmt.Errorf("update status: %w", err)
			}
		}

		if err := s.repo.UpdateProfile(ctx, id, displayName, username, status); err != nil {
			return fmt.Errorf("failed update user: %w", err)
		}

		updatedUser, err = s.repo.GetOwnProfile(ctx, id)
		if err != nil {
			return fmt.Errorf("get updated user: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}
