package fake

import (
	"context"
	"fmt"
	"time"

	"github.com/ilyaotinov/osync/internal/file"
)

type File struct {
	ModifyData time.Time
	NameData   string
	MD5Data    string
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
	Files           map[string]file.File
	AlwaysReturnErr bool
}

func (f *Filesystem) IsFileExists(ctx context.Context, path string) (bool, error) {
	if f.AlwaysReturnErr {
		return false, fmt.Errorf("failed check file for existence")
	}

	if _, ok := f.Files[path]; !ok {
		return false, nil
	}

	return true, nil
}

func (f *Filesystem) GetResource(ctx context.Context, path string) (file.File, error) {
	if f.AlwaysReturnErr {
		return nil, fmt.Errorf("internal error")
	}

	fInfo, ok := f.Files[path]
	if !ok {
		return nil, fmt.Errorf("f not found")
	}

	return fInfo, nil
}
