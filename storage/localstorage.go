package storage

import (
	"douyin/config"
	"fmt"
	"io"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
	"sync"
)

type LocalStorage struct {
}

func (LocalStorage) GetURL(path string) string {
	return fmt.Sprintf(
		"%s/%s",
		strings.TrimSuffix(config.Conf.StorageConfig.Local.BaseURL, "/"),
		path2.Join("static", path),
	)
}

func (LocalStorage) Upload(path string, reader io.Reader) error {
	path = filepath.Join(config.Conf.StorageConfig.Local.Path, path)

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

func (LocalStorage) Delete(path ...string) ([]string, error) {
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
