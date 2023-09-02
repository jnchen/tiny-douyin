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
	Delete(path ...string) error
}

var local = &Local{}
var oss *OSS

func GetStorage() Storage {
	if config.Conf.Storage.OSS.Enable {
		return oss
	}
	return local
}

func init() {
	if config.Conf.Storage.OSS.Enable {
		log.Println("初始化OSS")
		oss = &OSS{}
		var err error

		if oss.ossClient, err = aliyunoss.New(
			config.Conf.Storage.OSS.Endpoint,
			config.Conf.Storage.OSS.AccessKeyID,
			config.Conf.Storage.OSS.AccessKeySecret,
		); err != nil {
			log.Panicln("初始化OSSClient失败", err)
		}
		if oss.ossBucket, err = oss.ossClient.Bucket(
			config.Conf.Storage.OSS.BucketName,
		); err != nil {
			log.Panicln("初始化OSSBucket失败", err)
		}
	}
}
