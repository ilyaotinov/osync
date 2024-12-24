package yclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ilyaotinov/osync/internal/file"
)

type YandexClient struct {
	c       *http.Client
	baseURL string
	token   string
}

type GetResourceResponse struct {
	MD5      string    `json:"md5"`
	Modified time.Time `json:"modified"`
}

func New(client *http.Client, baseURL string, token string) *YandexClient {
	return &YandexClient{
		c:       client,
		baseURL: baseURL,
		token:   token,
	}
}

func (y *YandexClient) IsFileExists(ctx context.Context, path string) (bool, error) {
	fullURL, err := y.buildURL("/v1/disk/resources", map[string]string{
		"path": path,
	})
	if err != nil {
		return false, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create GET resource request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+y.token)

	resp, err := y.c.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to perform request to Yandex API: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body: ", "error", err.Error())
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

func (y *YandexClient) GetResource(ctx context.Context, path string) (*file.File, error) {
	fullURL, err := y.buildURL("/v1/disk/resources", map[string]string{
		"path": path,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET resource request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+y.token)

	resp, err := y.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET resource request: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed close response body: ", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	resource := &GetResourceResponse{}
	if err = json.NewDecoder(resp.Body).Decode(resource); err != nil {
		return nil, fmt.Errorf("failed to decode GET resource body: %w", err)
	}

	res := &file.File{}

	return res.SetModify(resource.Modified).SetMD5(resource.MD5), nil
}

func (y *YandexClient) buildURL(endpoint string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(y.baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	u.Path, err = url.JoinPath(u.Path, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed build url: %w", err)
	}

	query := u.Query()
	for key, value := range queryParams {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}
