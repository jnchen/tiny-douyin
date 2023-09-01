package storage

import (
	"douyin/config"
	"fmt"
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
	return filepath.Join(config.Conf.StorageConfig.Local.Path, path)
}

func (*Local) GetURL(path string) string {
	scheme, baseURLWithoutScheme, found := strings.Cut(
		config.Conf.StorageConfig.Local.BaseURL,
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

func (*Local) Upload(path string, reader io.Reader) error {
	path = filepath.Join(config.Conf.StorageConfig.Local.Path, path)
	return SaveAsLocalFile(path, reader)
}

func (*Local) Delete(path ...string) ([]string, error) {
	var totalErr error
	success := make([]string, len(path))
	var wg sync.WaitGroup

	for i, p := range path {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			fullPath := filepath.Join(config.Conf.StorageConfig.Local.Path, p)
			if err := os.Remove(fullPath); err != nil {
				totalErr = fmt.Errorf("%v\n%v", totalErr, err)
			} else {
				success[i] = p
			}
		}(i, p)
	}
	wg.Wait()

	var deletedFiles []string
	for _, p := range success {
		if p != "" {
			deletedFiles = append(deletedFiles, p)
		}
	}

	return deletedFiles, totalErr
}
