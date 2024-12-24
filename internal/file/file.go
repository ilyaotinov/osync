package file

import "time"

type File struct {
	modify time.Time
	md5    string
}

func (f *File) SetMD5(md5 string) *File {
	f.md5 = md5

	return f
}

func (f *File) SetModify(modify time.Time) *File {
	f.modify = modify

	return f
}

func (f *File) GetModify() time.Time {
	return f.modify
}

func (f *File) GetMD5() string {
	return f.md5
}

type ModifyInfo struct {
	ModifyDate time.Time
	Hash       string
}
