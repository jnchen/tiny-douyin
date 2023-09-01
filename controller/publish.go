package controller

import (
	"bytes"
	"douyin/model"
	"douyin/service"
	"douyin/storage"
	"douyin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func handleUploading(
	req *model.PublishActionRequest,
	videoFilePath,
	coverFilePath string,
	store storage.Storage,
	result chan<- error,
) {
	defer close(result)

	// TODO: 失败后？
	file, err := req.Data.Open()
	if nil != err {
		log.Println("打开视频文件失败：", err)
		result <- err
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	if err = store.Upload(videoFilePath, file); err != nil {
		log.Println("存储视频文件失败：", err)
		result <- err
		return
	}

	var tempVideoFilePath string
	switch store.(type) {
	case *storage.Local:
		tempVideoFilePath = storage.GetLocalPath(videoFilePath)
	case *storage.OSS:
		tempVideoFilePath = filepath.Join(
			os.TempDir(),
			"douyin_tmp_"+videoFilePath,
		)
		log.Println("临时视频文件路径：", tempVideoFilePath)
		_, err = file.Seek(0, 0)
		if err != nil {
			log.Println("重置视频文件指针失败：", err)
			result <- err
			return
		}
		if err = storage.SaveAsLocalFile(tempVideoFilePath, file); err != nil {
			log.Println("保存临时视频文件失败：", err)
			result <- err
			return
		}
	}

	imgBytes, err := util.ReadSingleFrameAsBytes(tempVideoFilePath, 1)
	if nil != err {
		log.Println("获取视频封面失败：", err)
		result <- err
		return
	}

	if err = store.Upload(
		coverFilePath,
		bytes.NewReader(imgBytes),
	); nil != err {
		log.Println("保存视频封面失败：", err)
		result <- err
		return
	}

	result <- nil
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	var req model.PublishActionRequest
	if err := c.ShouldBind(&req); nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	user := util.GetUser(c)
	uuid := util.UUID()
	userDir := fmt.Sprintf(
		"videos/%s/%s",
		fmt.Sprintf("%d", user.Id),
		time.Now().Format("2006-01-02"),
	)
	videoFileName := uuid + filepath.Ext(req.Data.Filename)
	videoFilePath := filepath.ToSlash(filepath.Join(userDir, videoFileName))
	coverFileName := uuid + ".jpg"
	coverFilePath := filepath.ToSlash(filepath.Join(userDir, coverFileName))
	store := storage.GetStorage()
	playUrl := store.GetURL(videoFilePath)
	coverUrl := store.GetURL(coverFilePath)

	resultUploading := make(chan error)
	go handleUploading(
		&req,
		videoFilePath,
		coverFilePath,
		store,
		resultUploading,
	)

	// 插入视频信息
	if _, err := service.PublishAction(
		user.Id,
		playUrl,
		coverUrl,
		req.Title,
	); nil != err {
		if err := <-resultUploading; nil == err {
			deleted, err := store.Delete(
				videoFilePath,
				coverFilePath,
			)
			if nil != err {
				log.Println("删除文件失败：", err)
			}

			for _, path := range deleted {
				log.Println("删除文件成功：", path)
			}
		}
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  videoFileName + "上传成功！",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var req model.PublishListRequest
	if err := c.ShouldBind(&req); nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoPublishListDAO, err := service.PublishList(req.UserId)
	if nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoPublishList := make([]model.Video, len(videoPublishListDAO))
	for i, video := range videoPublishListDAO {
		videoPublishList[i] = *video.ToModel()
	}

	c.JSON(http.StatusOK, model.PublishListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "获取视频列表成功！",
		},
		VideoList: videoPublishList,
	})
}
