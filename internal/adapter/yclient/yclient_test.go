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
const (
	existedPath       = "/test_private_osync/item_1.docx"
	existedFolderPath = "/test_private_osync"
)

func TestYandexClient_IsFileExistsIntegration(t *testing.T) {
	token := getToken(t)

	client := &http.Client{}
	yClient := New(client, baseYandexAPIURL, token)

	t.Run("file exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		got, err := yClient.IsFileExists(ctx, existedPath)

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

func TestYandexClient_GetResource_File(t *testing.T) {
	c := &http.Client{}
	token := getToken(t)
	yClient := New(c, baseYandexAPIURL, token)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	expectedModify, err := time.Parse(time.DateTime, "2020-01-01 00:00:00")
	require.NoError(t, err)

	got, err := yClient.GetResource(ctx, existedPath)

	require.NoErrorf(t, err, "expected file to be found and dont have error")
	assert.Truef(t, got.Modify().After(expectedModify), "expected file be modified after 2020 year")
	assert.Truef(t, len(got.MD5()) > 0, "expected to file hash not empty")
	assert.Falsef(t, got.IsDIR(), "expect file is dir was false")
}

func TestYandexClient_GetResource_Dir(t *testing.T) {
	c := &http.Client{}
	token := getToken(t)
	yClient := New(c, baseYandexAPIURL, token)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	expectedModify, err := time.Parse(time.DateTime, "2020-01-01 00:00:00")
	require.NoError(t, err)

	got, err := yClient.GetResource(ctx, existedFolderPath)

	require.NoErrorf(t, err, "expected file to be found and dont have error")
	assert.Truef(t, got.Modify().After(expectedModify), "expected folder be modified after 2020 year")
	assert.Truef(t, len(got.MD5()) == 0, "expected to folder hash be empty")
	assert.Truef(t, got.IsDIR(), "expect file is dir")
}

func getToken(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test for yandex client skipped...")
	}

	token := os.Getenv("Y_DISK_TOKEN")
	if len(token) == 0 {
		t.Skip("token for yandex disk not set. skipping...")
	}

	return token
}
