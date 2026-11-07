package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"uuid"
)

type ProfileService struct {
	repo usersRepo
}

type usersRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, displayName, username, status *string) error
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
}

func NewProfileService(repo usersRepo) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}
