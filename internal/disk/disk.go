package disk

import (
	"context"
	"fmt"

	"github.com/ilyaotinov/osync/internal/file"
)

type Filesystem interface {
	IsFileExists(ctx context.Context, path string) (bool, error)
	GetResource(ctx context.Context, path string) (*file.File, error)
}

type Disk struct {
	filesystem Filesystem
}

func New(filesystem Filesystem) *Disk {
	return &Disk{
		filesystem: filesystem,
	}
}

func (d *Disk) IsFileExists(ctx context.Context, path string) (bool, error) {
	if ctx == nil {
		return false, fmt.Errorf("ctx cannot be nil")
	}

	if len(path) == 0 {
		return false, fmt.Errorf("path cannot be empty")
	}

	res, err := d.filesystem.IsFileExists(ctx, path)
	if err != nil {
		return false, fmt.Errorf("failed check file for existence")
	}

	return res, nil
}

func (d *Disk) GetFileModificationInfo(ctx context.Context, path string) (file.ModifyInfo, error) {
	if ctx == nil {
		return file.ModifyInfo{}, fmt.Errorf("ctx cannot be nil")
	}

	if len(path) == 0 {
		return file.ModifyInfo{}, fmt.Errorf("path cannot be empty")
	}

	f, err := d.filesystem.GetResource(ctx, path)
	if err != nil {
		return file.ModifyInfo{}, fmt.Errorf("failed get file info from filesystem: %w", err)
	}

	return file.ModifyInfo{
		ModifyDate: f.GetModify(),
		Hash:       f.GetMD5(),
	}, nil
}
