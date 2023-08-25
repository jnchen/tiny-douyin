package controller

import (
	"bytes"
	"douyin/model"
	"douyin/service"
	"douyin/storage/oss"
	"douyin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

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

	file, err := req.Data.Open()
	if nil != err {
		log.Println("打开文件失败：", err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Panicf("关闭文件失败：%s\n", err)
		}
	}(file)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); nil != err {
		log.Println("读取文件失败：", err)
		return
	}

	data := buf.Bytes()
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
	playUrl := oss.GetURL(videoFilePath)
	coverUrl := oss.GetURL(coverFilePath)

	resultUploading := make(chan error)
	go func() {
		// TODO: 失败重传
		if err := oss.Upload(videoFilePath, data); err != nil {
			log.Println("存储视频文件失败：", err)
			resultUploading <- err
			return
		}
		imgBytes, err := util.ReadSingleFrameAsBytes(playUrl, 1)
		if nil != err {
			log.Println("获取视频封面失败：", err)
			resultUploading <- err
		} else {
			if err := oss.Upload(coverFilePath, imgBytes); nil != err {
				log.Println("保存视频封面失败：", err)
				resultUploading <- err
			}
		}
		resultUploading <- nil
	}()

	// 插入视频信息
	if _, err := service.VideoPublish(
		user.Id,
		playUrl,
		coverUrl,
		req.Title,
	); nil != err {
		if err := <-resultUploading; nil == err {
			if res, err := oss.Delete(videoFilePath, coverFilePath); nil != err {
				log.Println("删除文件失败：", err)
				log.Println("失败项目：")
				for _, obj := range res.DeletedObjects {
					log.Println("\t", obj)
				}
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

	videoPublishListDAO, err := service.VideoPublishList(req.UserId)
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

	if user := util.GetUser(c); nil != user {
		for i, video := range videoPublishList {
			videoPublishList[i].IsFavorite, _ = service.FavoriteCheck(user.Id, video.Id)
		}
	}

	c.JSON(http.StatusOK, model.PublishListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "获取视频列表成功！",
		},
		VideoList: videoPublishList,
	})
}
