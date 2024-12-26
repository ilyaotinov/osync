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
		expectedModify.After(resource.GetModify()),
		"expected modify time to be after actual modify time",
	)
	assert.Equal(t, expectedHash, resource.GetMD5())
}
