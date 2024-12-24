package yclient

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseYandexAPIURL = "https://cloud-api.yandex.net"

func TestYandexClient_IsFileExistsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test for yandex client skipped...")
	}

	token := os.Getenv("Y_DISK_TOKEN")
	if len(token) == 0 {
		t.Skip("token for yandex disk not set. skipping...")
	}

	client := &http.Client{}
	yClient := New(client, baseYandexAPIURL, token)

	t.Run("file exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		got, err := yClient.IsFileExists(ctx, "/test_private_osync/item_1.docx")

		require.NoErrorf(t, err, "unexpected error from yclient")
		assert.Truef(t, got, "expected file be found on yandex disk")
	})

	t.Run("file not exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		got, err := yClient.IsFileExists(ctx, "/test_private_osync/item_1_not_exists.docx")

		require.NoErrorf(t, err, "unexpected error from yclient")
		assert.Falsef(t, got, "expected file not be found on yandex disk")
	})

	t.Run("folder exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		got, err := yClient.IsFileExists(ctx, "/test_private_osync")

		require.NoErrorf(t, err, "unexpected error from yclient")
		assert.Truef(t, got, "expected folder be found on yandex disk")
	})

	t.Run("error becouse of timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		_, err := yClient.IsFileExists(ctx, "/test_private_osync")

		require.Errorf(t, err, "expected error from yclient becouse of context timeout")
	})
}
