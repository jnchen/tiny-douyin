package storage

import (
	"fmt"
	"strings"
)

type UploadingError struct {
	err  error
	path string
}

type DeletingError struct {
	msg  *string
	errs []error
}

func (e UploadingError) Error() string {
	return fmt.Sprintf("上传 %s 失败：%s", e.path, e.err.Error())
}

func (e UploadingError) Unwrap() error {
	return e.err
}

func (e DeletingError) Error() string {
	if e.msg != nil {
		return *e.msg
	}

	var res strings.Builder
	for _, err := range e.errs {
		res.WriteString(err.Error())
		res.WriteRune('\n')
	}

	e.msg = new(string)
	*e.msg = res.String()
	return *e.msg
}

func (e DeletingError) Unwrap() []error {
	return e.errs
}
