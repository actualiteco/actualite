package profile_service

import (
	"actualiteco/actualite/internal/core/domain"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"uuid"
)

var avatarExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func (s *ProfileService) UploadAvatar(ctx context.Context, id uuid.UUID, data []byte) (string, error) {
	if len(data) == 0 || len(data) > domain.MaxAvatarSize {
		return "", fmt.Errorf("invalid avatar size=%d: %w", len(data), domain.ErrInvalidRequest)
	}

	contentType := http.DetectContentType(data)
	ext, ok := avatarExtensions[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported avatar content type=%s: %w", contentType, domain.ErrInvalidRequest)
	}

	avatarURL, err := s.storage.Upload(ctx, "avatar"+ext, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("upload avatar to storage: %w", err)
	}

	if err := s.repo.UpdateAvatar(ctx, id, avatarURL); err != nil {
		return "", fmt.Errorf("update avatar url: %w", err)
	}

	return avatarURL, nil
}
