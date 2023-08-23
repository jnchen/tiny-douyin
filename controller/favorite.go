package controller

import (
	"douyin/model"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FavoriteAction(c *gin.Context) {
	var req model.FavoriteActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	user, exist := service.CheckLogin(req.Token)
	if !exist {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

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

	_, exist := service.CheckLogin(req.Token)
	if !exist {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

	videoDAOList, err := service.FavoriteList(req.UserId)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoList := make([]model.Video, len(videoDAOList))
	for i, videoDAO := range videoDAOList {
		videoList[i] = *videoDAO.ToModel()
		videoList[i].IsFavorite = true
	}

	c.JSON(http.StatusOK, model.FavoriteListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
