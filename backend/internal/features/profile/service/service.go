package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"context"
	"io"
	"uuid"
)

type ProfileService struct {
	repo      usersRepo
	storage   fileStorage
	txManager txManager
}

type fileStorage interface {
	Upload(ctx context.Context, filename string, file io.Reader) (string, error)
}

type txManager interface {
	DoWrite(ctx context.Context, fn func(context.Context) error) error
}

type usersRepo interface {
	GetOwnProfile(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByID(ctx context.Context, id, targetID uuid.UUID) (domain.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, displayName, username, status *string) error
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsUsernameTaken(ctx context.Context, username string) (bool, error)
	UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) error
}

func NewProfileService(repo usersRepo, storage fileStorage, txManager txManager) *ProfileService {
	return &ProfileService{
		repo:      repo,
		storage:   storage,
		txManager: txManager,
	}
}
