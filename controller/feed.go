package controller

import (
	"douyin/model"
	"douyin/service"
	"douyin/util"
	"github.com/gin-gonic/gin"
	"github.com/u2takey/go-utils/integer"
	"net/http"
	"time"
)

const FeedLimit = 10

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var req model.FeedRequest
	if err := c.ShouldBindQuery(&req); nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	latestTime := time.UnixMilli(req.LatestTime)
	if 0 == req.LatestTime {
		latestTime = time.Now()
	}
	// log.Println("请求时间", latestTime.UnixMilli())

	var userId int64 = -1
	user := util.GetUser(c)
	if user != nil {
		userId = user.Id
	}

	videoListDAO, err := service.FeedList(userId, latestTime, FeedLimit)
	if nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	n := len(videoListDAO)
	videoList := make([]model.Video, n)
	var nextTime int64 = 0
	if n > 0 {
		nextTime = videoListDAO[0].CreatedAt.UnixMilli()
	}
	for i, video := range videoListDAO {
		// 视频信息是按照创建时间倒序排列的，虽然可以直接获取最后一个视频的创建时间，
		// 但是为了不依赖于数据库的实现，这里还是遍历一遍
		nextTime = integer.Int64Min(nextTime, video.CreatedAt.UnixMilli())
		videoList[i] = *video.ToModel()
		// log.Println("视频", video.ID, "创建时间", video.CreatedAt.UnixMilli())
	}
	// log.Println("下次请求时间", time.UnixMilli(nextTime).UnixMilli())

	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
