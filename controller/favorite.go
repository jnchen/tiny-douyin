package controller

import (
	"douyin/model"
	"douyin/service"
	"douyin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FavoriteAction(c *gin.Context) {
	var req model.FavoriteActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	user := util.GetUser(c)
	if req.ActionType == 1 {
		if err := service.FavoriteAction(user.Id, req.VideoId); err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	} else if req.ActionType == 2 {
		if err := service.FavoriteDelete(user.Id, req.VideoId); err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var req model.FavoriteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoDAOList, err := service.FavoriteList(req.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoList := make([]model.Video, len(videoDAOList))
	for i, videoDAO := range videoDAOList {
		videoList[i] = *videoDAO.ToModel()
	}

	c.JSON(http.StatusOK, model.FavoriteListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
