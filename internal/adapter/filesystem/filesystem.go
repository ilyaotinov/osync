package filesystem

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ilyaotinov/osync/internal/file"
)

type File struct {
	ModifyData time.Time
	MD5Data    string
	NameData   string
	IsDIRData  bool
}

func (f File) Name() string {
	return f.NameData
}

func (f File) Modify() time.Time {
	return f.ModifyData
}

func (f File) MD5() string {
	return f.MD5Data
}

func (f File) IsDIR() bool {
	return f.IsDIRData
}

type Filesystem struct {
}

func New() *Filesystem {
	return &Filesystem{}
}

func (f Filesystem) IsFileExists(ctx context.Context, path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("error check file for existence: %w", err)
	}
}

func (f Filesystem) GetResource(ctx context.Context, path string) (file.File, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	var hash string
	if !fInfo.IsDir() {
		hash, err = getFileMD5(path)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate hash of file: %w", err)
		}
	}

	result := File{
		ModifyData: fInfo.ModTime(),
		MD5Data:    hash,
		IsDIRData:  fInfo.IsDir(),
		NameData:   fInfo.Name(),
	}

	return result, nil
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
