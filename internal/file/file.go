package file

import "time"

type File interface {
	Modify() time.Time
	MD5() string
	IsDIR() bool
}
