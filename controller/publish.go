package controller

import (
	"douyin/config"
	"douyin/model"
	"douyin/service"
	"douyin/util"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	data := req.Data
	user := util.GetUser(c)
	uuid := util.UUID()
	userDir := filepath.Join("./public/videos", fmt.Sprintf("%d", user.Id))
	if err := os.MkdirAll(userDir, os.ModePerm); nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoFileName := fmt.Sprintf("%s%s", uuid, filepath.Ext(data.Filename))
	videoFilePath := filepath.Join(userDir, videoFileName)
	if err := c.SaveUploadedFile(data, videoFilePath); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	playUrl := fmt.Sprintf("%s%s",
		config.Conf.BaseURL,
		path.Join("/static/videos/", fmt.Sprintf("%d", user.Id), videoFileName),
	)

	coverFileName := fmt.Sprintf("%s%s", uuid, ".jpg")
	coverFilePath := filepath.Join(userDir, coverFileName)
	coverUrl := "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	if img, err := util.ReadVideoSingleFrame(videoFilePath, 0); nil != err {
		log.Println("获取视频封面失败：", err)
	} else {
		if err := imaging.Save(img, coverFilePath); nil != err {
			log.Println("保存视频封面失败：", err)
		} else {
			coverUrl = fmt.Sprintf("%s%s",
				config.Conf.BaseURL,
				path.Join("/static/videos/", fmt.Sprintf("%d", user.Id), coverFileName),
			)
		}
	}

	// 插入视频信息
	if _, err := service.VideoPublish(
		user.Id,
		playUrl,
		coverUrl,
		req.Title,
	); nil != err {
		// 删除视频文件
		if err := os.Remove(videoFilePath); nil != err {
			log.Printf("删除视频文件%s失败：%s\n", videoFilePath, err)
		}
		// 删除封面文件
		if err := os.Remove(coverFilePath); nil != err {
			log.Printf("删除封面文件%s失败：%s\n", coverFilePath, err)
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
