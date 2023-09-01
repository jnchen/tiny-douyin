package controller

import (
	"douyin/model"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"github.com/u2takey/go-utils/integer"
	"net/http"
	"time"
)

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

	const limit = 10
	latestTime := time.UnixMilli(req.LatestTime)
	if 0 == req.LatestTime {
		latestTime = time.Now()
	}

	var userId int64 = -1
	user, err := service.CheckLogin(req.Token)
	if nil == err {
		userId = user.Id
	}

	videoListDAO, err := service.FeedList(userId, latestTime, limit)
	if nil != err {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	n := len(videoListDAO)
	videoList := make([]model.Video, n)
	var nextTime int64 = 0 // 单位：秒
	if n > 0 {
		nextTime = videoListDAO[0].CreatedAt.Unix()
	}
	for i, video := range videoListDAO {
		// 视频信息是按照创建时间倒序排列的，虽然可以直接获取最后一个视频的创建时间，
		// 但是为了不依赖于数据库的实现，这里还是遍历一遍
		nextTime = integer.Int64Min(nextTime, video.CreatedAt.Unix())
		videoList[i] = *video.ToModel()
	}
	// fmt.Println("下次请求时间", time.Unix(nextTime, 0).Format("2006-01-02 15:04:05"))

	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime * 1000,
	})
}
