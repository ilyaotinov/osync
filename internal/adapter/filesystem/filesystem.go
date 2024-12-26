package filesystem

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/ilyaotinov/osync/internal/file"
)

type Filesystem struct {
}

func New() *Filesystem {
	return &Filesystem{}
}

func (f Filesystem) IsFileExists(ctx context.Context, path string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (f Filesystem) GetResource(ctx context.Context, path string) (*file.File, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	hash, err := getFileMD5(path)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash of file: %w", err)
	}

	result := &file.File{}

	return result.SetModify(fInfo.ModTime()).SetMD5(hash), nil
}

func getFileMD5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open f: %w", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			slog.Error("unable to close file: ", "error", err)
		}
	}()

	hash := md5.New()

	if _, err = io.Copy(hash, f); err != nil {
		return "", fmt.Errorf("unable to compute hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
