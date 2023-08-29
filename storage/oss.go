package storage

import (
	"douyin/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

type OSSStorage struct {
}

func (OSSStorage) GetURL(path string) string {
	return fmt.Sprintf(
		"https://%s.%s/%s",
		config.Conf.StorageConfig.OSS.BucketName,
		config.Conf.StorageConfig.OSS.Endpoint,
		path,
	)
}

func (OSSStorage) Upload(path string, reader io.Reader) error {
	return ossBucket.PutObject(
		path,
		reader,
		oss.ObjectACL(oss.ACLPublicRead),
	)
}

func (OSSStorage) Delete(path ...string) ([]string, error) {
	res, err := ossBucket.DeleteObjects(path)
	return res.DeletedObjects, err
}
