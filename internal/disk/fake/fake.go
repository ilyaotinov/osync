package fake

import (
	"context"
	"fmt"

	"github.com/ilyaotinov/osync/internal/file"
)

type FakeFilesystem struct {
	Files           map[string]*file.File
	AlwaysReturnErr bool
}

func (f *FakeFilesystem) IsFileExists(ctx context.Context, path string) (bool, error) {
	if f.AlwaysReturnErr {
		return false, fmt.Errorf("failed check file for existence")
	}

	if _, ok := f.Files[path]; !ok {
		return false, nil
	}

	return true, nil
}

func (f *FakeFilesystem) GetResource(ctx context.Context, path string) (*file.File, error) {
	if f.AlwaysReturnErr {
		return nil, fmt.Errorf("internal error")
	}

	file, ok := f.Files[path]
	if !ok {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}
