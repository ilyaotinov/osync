package disk

import (
	"context"
	"testing"
	"time"

	"github.com/ilyaotinov/osync/internal/disk/fake"
	"github.com/ilyaotinov/osync/internal/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisk_IsFileExists(t *testing.T) {
	type args struct {
		ctxFunc func() context.Context
		path    string
	}

	type expect struct {
		val bool
		err bool
	}

	tests := []struct {
		name       string
		filesystem *fake.FakeFilesystem
		args       args
		expect     expect
	}{
		{
			name: "file found case",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"/path": {},
				},
				AlwaysReturnErr: false,
			},
			args:   args{ctxFunc: context.Background, path: "/path"},
			expect: expect{val: true, err: false},
		},
		{
			name: "file not found case",
			filesystem: &fake.FakeFilesystem{
				Files:           map[string]*file.File{},
				AlwaysReturnErr: false,
			},
			args: args{
				ctxFunc: context.Background,
				path:    "/path",
			},
			expect: expect{
				val: false,
				err: false,
			},
		},
		{
			name: "null ctx given",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"/path": {},
				},
				AlwaysReturnErr: false,
			},
			args: args{
				ctxFunc: func() context.Context {
					return nil
				},
				path: "/path",
			},
			expect: expect{
				val: false,
				err: true,
			},
		},
		{
			name:       "empty path given",
			filesystem: &fake.FakeFilesystem{},
			args: args{
				ctxFunc: context.Background,
				path:    "",
			},
			expect: expect{
				val: false,
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disk := New(tt.filesystem)
			got, err := disk.IsFileExists(tt.args.ctxFunc(), tt.args.path)
			if tt.expect.err {
				require.Errorf(t, err, "expect get error from file exists fucnction")
				return
			} else {
				require.NoErrorf(t, err, "expect no error from file exists function")
			}
			assert.Equalf(t, tt.expect.val, got, "unexpected file existence check result")
		})
	}
}

func TestDisk_GetFileModificationInfo(t *testing.T) {
	type expect struct {
		hash string
		mod  time.Time
		err  bool
	}

	type arg struct {
		ctxFunc func() context.Context
		path    string
	}

	tests := []struct {
		name       string
		filesystem *fake.FakeFilesystem
		arg        arg
		expect     expect
	}{
		{
			name: "success case",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"/path": func() *file.File {
						mod, _ := time.Parse(time.DateTime, "2024-01-02 00:00:00")
						file := &file.File{}
						return file.SetMD5("hash-expectedd").SetModify(mod)
					}(),
				},
				AlwaysReturnErr: false,
			},
			arg: arg{ctxFunc: context.Background, path: "/path"},
			expect: expect{
				hash: "hash-expectedd",
				mod: func() time.Time {
					res, _ := time.Parse(time.DateTime, "2024-01-02 00:00:00")

					return res
				}(),
				err: false,
			},
		},
		{
			name: "null ctx given",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"/path": func() *file.File {
						f := &file.File{}
						return f.SetMD5("test").SetModify(time.Now())
					}(),
				},
				AlwaysReturnErr: false,
			},
			arg: arg{
				ctxFunc: func() context.Context {
					return nil
				},
				path: "/path",
			},
			expect: expect{
				hash: "",
				mod:  time.Time{},
				err:  true,
			},
		},
		{
			name: "empty path",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"": func() *file.File {
						f := &file.File{}
						return f.SetModify(time.Now()).SetMD5("test-hash")
					}(),
				},
				AlwaysReturnErr: false,
			},
			arg: arg{
				ctxFunc: context.Background,
				path:    "",
			},
			expect: expect{
				err: true,
			},
		},
		{
			name: "internal file storage error",
			filesystem: &fake.FakeFilesystem{
				Files: map[string]*file.File{
					"/path": func() *file.File {
						f := &file.File{}
						return f.SetMD5("hash").SetModify(time.Now())
					}(),
				},
				AlwaysReturnErr: true,
			},
			arg: arg{
				ctxFunc: context.Background,
				path:    "/path",
			},
			expect: expect{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New(tt.filesystem)

			got, err := d.GetFileModificationInfo(tt.arg.ctxFunc(), tt.arg.path)

			if tt.expect.err {
				require.Error(t, err)

				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expect.mod, got.ModifyDate)
			assert.Equal(t, tt.expect.hash, got.Hash)
		})
	}
}
