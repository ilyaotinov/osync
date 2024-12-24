package yclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type YandexClient struct {
	c       *http.Client
	baseURL string
	token   string
}

type GetResourceResponse struct {
	Modify time.Time
	MD5    string
}

func New(client *http.Client, baseURL string, token string) *YandexClient {
	return &YandexClient{
		c:       client,
		baseURL: baseURL,
		token:   token,
	}
}

func (y *YandexClient) IsFileExists(ctx context.Context, path string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, y.baseURL+"/v1/disk/resources?path="+path, nil)
	if err != nil {
		return false, fmt.Errorf("failed create get resource request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+y.token)

	resp, err := y.c.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed perform request for check existence to yandex api: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed close response body: ", "error", err.Error())
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %w", err)
	}

	return true, nil
}
