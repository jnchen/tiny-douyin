package storage

import (
	"douyin/config"
	"fmt"
	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type Storage interface {
	GetURL(path string) string
	Upload(path string, reader io.Reader) error
	Delete(path ...string) error
}

var local *Local
var oss *OSS
var store Storage

func GetLocalStorage() *Local {
	return local
}

func GetStorage() Storage {
	return store
}

func Init(storageConfig *config.Storage) error {
	var err error

	local = &Local{
		localRoot: storageConfig.Local.Path,
		baseURL:   storageConfig.Local.BaseURL,
	}

	if storageConfig.OSS.Enable {
		log.Println("初始化OSS")

		oss = &OSS{
			bucketName: storageConfig.OSS.BucketName,
			endpoint:   storageConfig.OSS.Endpoint,
		}
		if oss.ossClient, err = aliyunoss.New(
			storageConfig.OSS.Endpoint,
			storageConfig.OSS.AccessKeyID,
			storageConfig.OSS.AccessKeySecret,
		); err != nil {
			return fmt.Errorf("初始化 OSSClient 失败：%w", err)
		}
		if oss.ossBucket, err = oss.ossClient.Bucket(
			storageConfig.OSS.BucketName,
		); err != nil {
			return fmt.Errorf("初始化 OSSBucket 失败：%w", err)
		}
		store = oss
	} else {
		store = local
	}
	return err
}
