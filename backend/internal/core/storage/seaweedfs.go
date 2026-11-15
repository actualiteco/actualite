package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type SeaweedFSClient struct {
	masterURL    string
	publicURL    string
	volumeServer string
	client       *http.Client
}

type assignResponse struct {
	Fid string `json:"fid"`
}

func NewSeaweedFSClient(cfg Config) *SeaweedFSClient {
	return &SeaweedFSClient{
		masterURL:    strings.TrimRight(cfg.MasterURL, "/"),
		publicURL:    strings.TrimRight(cfg.PublicURL, "/"),
		volumeServer: cfg.VolumeServer,
		client:       &http.Client{Timeout: cfg.Timeout},
	}
}

func (s *SeaweedFSClient) Upload(ctx context.Context, filename string, file io.Reader) (string, error) {
	fid, err := s.assignFileID(ctx)
	if err != nil {
		return "", fmt.Errorf("assign file id: %w", err)
	}

	if err := s.uploadToVolume(ctx, fid, filename, file); err != nil {
		return "", fmt.Errorf("upload to volume: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.publicURL, fid), nil
}

func (s *SeaweedFSClient) assignFileID(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.masterURL+"/dir/assign", nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request assign: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("assign failed with status: %d", resp.StatusCode)
	}

	var assign assignResponse
	if err := json.NewDecoder(resp.Body).Decode(&assign); err != nil {
		return "", fmt.Errorf("decode assign response: %w", err)
	}

	if assign.Fid == "" {
		return "", fmt.Errorf("assign response has empty fid")
	}

	return assign.Fid, nil
}

func (s *SeaweedFSClient) uploadToVolume(ctx context.Context, fid, filename string, file io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}

	uploadURL := fmt.Sprintf("http://%s/%s", s.volumeServer, fid)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	return nil
}
