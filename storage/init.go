package storage

import (
	"douyin/config"
	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type Storage interface {
	GetURL(path string) string
	Upload(path string, reader io.Reader) error
	Delete(path ...string) (deleted []string, err error)
}

var local = &Local{}
var oss *OSS

func GetStorage() Storage {
	if config.Conf.StorageConfig.OSS.Enable {
		return oss
	}
	return local
}

func init() {
	if config.Conf.StorageConfig.OSS.Enable {
		log.Println("初始化OSS")
		oss = &OSS{}
		var err error

		if oss.ossClient, err = aliyunoss.New(
			config.Conf.StorageConfig.OSS.Endpoint,
			config.Conf.StorageConfig.OSS.AccessKeyID,
			config.Conf.StorageConfig.OSS.AccessKeySecret,
		); err != nil {
			log.Panicln("初始化OSSClient失败", err)
		}
		if oss.ossBucket, err = oss.ossClient.Bucket(
			config.Conf.StorageConfig.OSS.BucketName,
		); err != nil {
			log.Panicln("初始化OSSBucket失败", err)
		}
	}
}
