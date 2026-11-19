package post_service

import (
	"actualiteco/actualite/internal/core/domain"
	"bytes"
	"context"
	"fmt"
	"net/http"
)

var imageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func (s *PostsService) UploadImage(ctx context.Context, data []byte) (string, error) {
	contentType := http.DetectContentType(data)
	ext, ok := imageExtensions[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported image content type=%s: %w", contentType, domain.ErrUnsupportedImageType)
	}

	imageURL, err := s.storage.Upload(ctx, "image"+ext, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("upload avatar to storage: %w", err)
	}

	return imageURL, nil
}
