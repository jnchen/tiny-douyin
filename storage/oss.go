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
		config.Conf.Storage.OSS.BucketName,
		config.Conf.Storage.OSS.Endpoint,
		path,
	)
}

func (o *OSS) Upload(path string, reader io.Reader) error {
	if err := o.ossBucket.PutObject(
		path,
		reader,
		aliyunoss.ObjectACL(aliyunoss.ACLPublicRead),
	); err != nil {
		return UploadingError{err, path}
	}

	return nil
}

func (o *OSS) Delete(path ...string) error {
	if len(path) == 0 {
		return nil
	}

	var deletingErr DeletingError
	res, err := o.ossBucket.DeleteObjects(path)
	if err != nil {
		deletingErr.errs = append(deletingErr.errs, err)
		return deletingErr
	}

	numDeleted := len(res.DeletedObjects)
	if numDeleted == len(path) {
		return nil
	}

	// 找出删除失败的对象
	deletedObjects := make(map[string]struct{}, numDeleted)
	for _, key := range res.DeletedObjects {
		deletedObjects[key] = struct{}{}
	}
	for _, p := range path {
		if _, ok := deletedObjects[p]; !ok {
			deletingErr.errs = append(
				deletingErr.errs,
				fmt.Errorf("删除对象 %s 失败", p),
			)
		}
	}
	return deletingErr
}
