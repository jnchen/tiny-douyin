package storage

import (
	"douyin/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type Storage interface {
	GetURL(path string) string
	Upload(path string, reader io.Reader) error
	Delete(path ...string) (deleted []string, err error)
}

var Impl Storage

var ossClient *oss.Client
var ossBucket *oss.Bucket

func init() {
	if config.Conf.StorageConfig.OSS.Enable {
		log.Println("初始化OSS")
		var err error
		ossClient, err = oss.New(
			config.Conf.StorageConfig.OSS.Endpoint,
			config.Conf.StorageConfig.OSS.AccessKeyID,
			config.Conf.StorageConfig.OSS.AccessKeySecret,
		)
		if err != nil {
			log.Panicln("初始化OSSClient失败", err)
		}
		ossBucket, err = ossClient.Bucket(config.Conf.StorageConfig.OSS.BucketName)
		if err != nil {
			log.Panicln("初始化OSSBucket失败", err)
		}
		Impl = OSSStorage{}
	} else {
		log.Println("初始化LocalStorage")
		Impl = LocalStorage{}
	}
}
