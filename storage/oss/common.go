package oss

import (
	"bytes"
	"douyin/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

func GetURL(path string) string {
	return fmt.Sprintf(
		"https://%s.%s/%s",
		config.Conf.BucketName,
		config.Conf.Endpoint,
		path,
	)
}

func Upload(path string, data []byte) error {
	err := ossBucket.PutObject(
		path,
		bytes.NewReader(data),
		oss.ObjectACL(oss.ACLPublicRead),
	)
	if err != nil {
		log.Println("上传数据失败", err)
		return err
	}
	return nil
}

func Delete(path ...string) (oss.DeleteObjectsResult, error) {
	res, err := ossBucket.DeleteObjects(path)
	if err != nil {
		log.Println("删除数据失败", err)
	}
	return res, err
}
