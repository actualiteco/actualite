package httprequest

import (
	"actualiteco/actualite/internal/core/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, domain.ErrInvalidRequest)
	}

	err := requestValidator.Struct(dest)
	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, domain.ErrInvalidRequest)
	}

	return nil
}

func GetIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("get user ip: %v: %w", err, domain.ErrInvalidRequest)
	}
	return ip, nil
}

func GetFileBytes(r *http.Request, key string, maxSize int64) ([]byte, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return nil, fmt.Errorf("upload exceeds limit: %w", domain.ErrRequestIsTooLarge)
		}

		return nil, fmt.Errorf("get file from form: %w: %w", err, domain.ErrInvalidRequest)
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	return readFileLimited(file, maxSize)
}

func readFileLimited(file io.Reader, maxSize int64) ([]byte, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("invalid max size=%d", maxSize)
	}

	data, err := io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return nil, fmt.Errorf("read file: %w: %w", err, domain.ErrInvalidRequest)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty file: %w", domain.ErrInvalidRequest)
	}

	if int64(len(data)) > maxSize {
		return nil, fmt.Errorf("file size=%d exceeds limit=%d: %w", len(data), maxSize, domain.ErrRequestIsTooLarge)
	}

	return data, nil
}
