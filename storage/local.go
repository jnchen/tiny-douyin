package storage

import (
	"douyin/util"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Local struct {
	localRoot string
	baseURL   string
}

func (l *Local) GetLocalPath(path string) string {
	return filepath.Join(l.localRoot, path)
}

func (l *Local) GetURL(path string) string {
	scheme, baseURLWithoutScheme, found := strings.Cut(
		l.baseURL,
		"://",
	)
	if !found {
		baseURLWithoutScheme = scheme
		scheme = "https"
	}

	u := url.URL{
		Scheme: scheme,
		Host:   baseURLWithoutScheme,
		Path: strings.Replace(
			filepath.ToSlash(l.GetLocalPath(path)),
			"public",
			"static",
			1,
		),
	}
	return u.String()
}

func (l *Local) Upload(path string, reader io.Reader) error {
	path = l.GetLocalPath(path)
	if err := util.SaveAsLocalFile(path, reader); err != nil {
		return UploadingError{err, path}
	}
	return nil
}

func (l *Local) Delete(path ...string) error {
	if len(path) == 0 {
		return nil
	}

	deletingErrs := util.NewConcurrentSlice[error]()
	var wg sync.WaitGroup

	for _, p := range path {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			fullPath := l.GetLocalPath(p)
			err := os.Remove(fullPath)
			if err == nil {
				return
			}
			deletingErrs.Append(err)
		}(p)
	}
	wg.Wait()

	if deletingErrs.Len() > 0 {
		return DeletingError{errs: deletingErrs.RawSlice()}
	}
	return nil
}
