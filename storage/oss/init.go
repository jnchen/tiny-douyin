package oss

import (
	"douyin/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

var ossClient *oss.Client
var ossBucket *oss.Bucket

func init() {
	log.Println("初始化OSS")
	var err error
	ossClient, err = oss.New(
		config.Conf.Endpoint,
		config.Conf.AccessKeyID,
		config.Conf.AccessKeySecret,
	)
	if err != nil {
		log.Panicln("初始化OSSClient失败", err)
	}
	ossBucket, err = ossClient.Bucket(config.Conf.BucketName)
	if err != nil {
		log.Panicln("初始化OSSBucket失败", err)
	}
}
