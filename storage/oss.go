package storage

import (
	"douyin/config"
	"fmt"
	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

type OSS struct {
	ossClient *aliyunoss.Client
	ossBucket *aliyunoss.Bucket
}

func (*OSS) GetURL(path string) string {
	return fmt.Sprintf(
		"https://%s.%s/%s",
		config.Conf.StorageConfig.OSS.BucketName,
		config.Conf.StorageConfig.OSS.Endpoint,
		path,
	)
}

func (o *OSS) Upload(path string, reader io.Reader) error {
	return o.ossBucket.PutObject(
		path,
		reader,
		aliyunoss.ObjectACL(aliyunoss.ACLPublicRead),
	)
}

func (o *OSS) Delete(path ...string) ([]string, error) {
	res, err := o.ossBucket.DeleteObjects(path)
	return res.DeletedObjects, err
}
