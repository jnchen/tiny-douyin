package storage

import (
	"douyin/config"
	"douyin/util"
	"io"
	"net/url"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
	"sync"
)

type Local struct {
}

func GetLocalPath(path string) string {
	return filepath.Join(config.Conf.Storage.Local.Path, path)
}

func (Local) GetURL(path string) string {
	scheme, baseURLWithoutScheme, found := strings.Cut(
		config.Conf.Storage.Local.BaseURL,
		"://",
	)
	if !found {
		baseURLWithoutScheme = scheme
		scheme = "https"
	}

	u := url.URL{
		Scheme: scheme,
		Host:   baseURLWithoutScheme,
		Path:   path2.Join("static", path),
	}
	return u.String()
}

func SaveAsLocalFile(path string, reader io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); nil != err {
		return err
	}

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	_, err = io.Copy(dst, reader)
	return err
}

func (Local) Upload(path string, reader io.Reader) error {
	path = filepath.Join(config.Conf.Storage.Local.Path, path)
	if err := SaveAsLocalFile(path, reader); err != nil {
		return UploadingError{err, path}
	}
	return nil
}

func (Local) Delete(path ...string) error {
	if len(path) == 0 {
		return nil
	}

	deletingErrs := util.NewConcurrentSlice[error]()
	var wg sync.WaitGroup

	for _, p := range path {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			fullPath := filepath.Join(config.Conf.Storage.Local.Path, p)
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
