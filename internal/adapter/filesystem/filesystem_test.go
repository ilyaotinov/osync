package filesystem

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilesystem_GetResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file system tests in short mode")
	}

	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer func() {
		if err = tmpFile.Close(); err != nil {
			t.Errorf("failed to close temporary file: %v", err)
		}
	}()
	_, err = tmpFile.Write([]byte("test data for hash calculation check"))
	require.NoError(t, err)

	t.Cleanup(func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("failed to remove temporary file: %v", err)
		}
	})

	expectedModify := time.Now().Add(time.Second * 1)
	expectedHash := "5c26ed997602442275b6639413926fa9"

	ctx := context.Background()
	path := tmpFile.Name()

	fsystem := New()

	resource, err := fsystem.GetResource(ctx, path)

	require.NoError(t, err)
	assert.Truef(
		t,
		expectedModify.After(resource.Modify()),
		"expected modify time to be after actual modify time",
	)
	assert.Equal(t, expectedHash, resource.MD5())
	assert.Falsef(t, resource.IsDIR(), "expect file has is dir false")
}

func TestFilesystem_GetResource_GetDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file system tests in short mode")
	}

	tempDir, err := os.MkdirTemp("", "test-*")
	require.NoError(t, err)

	ctx := context.Background()
	fsystem := New()

	got, err := fsystem.GetResource(ctx, tempDir)

	require.NoError(t, err)
	assert.Equal(t, got.MD5(), "")
	assert.Truef(t,
		got.Modify().After(time.Now().Add(time.Second*-1)), "expected modify time to be after actual modify time")
	assert.Truef(t, got.IsDIR(), "expect directory has is dir true")
}

func TestFilesystem_IsFileExists_FileExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file system tests in short mode")
	}
	tempFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer func() {
		err = tempFile.Close()
		require.NoError(t, err)
	}()
	t.Cleanup(func() {
		if err = os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temporary file: %v", err)
		}
	})

	ctx := context.Background()

	fsystem := New()

	got, err := fsystem.IsFileExists(ctx, tempFile.Name())

	require.NoError(t, err)
	assert.Truef(t, got, "expected file to be exists")
}

func TestFilesystem_IsFileExists_FileDoesNotExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file system tests in short mode")
	}

	ctx := context.Background()
	fsystem := New()

	got, err := fsystem.IsFileExists(ctx, "/path")

	require.NoError(t, err)
	assert.Falsef(t, got, "expected file to not exists")
}

func TestFilesystem_IsFileExists_DirectoryExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file system tests in short mode")
	}

	tempDir, err := os.MkdirTemp("", "test_dir_*")
	require.NoError(t, err)

	t.Cleanup(func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to remove temporary directory: %v", err)
		}
	})

	ctx := context.Background()
	fsystem := New()

	got, err := fsystem.IsFileExists(ctx, tempDir)

	require.NoError(t, err)
	assert.Truef(t, got, "expected file to be exists")
}
