package test

import (
	"douyin/controller"
	"douyin/model"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"path/filepath"
	"testing"
)

type BaseAPITestSuite struct {
	APITestSuite
}

func (s *BaseAPITestSuite) testUserAction() {
	e := s.newHTTPExpect(httpexpect.NewDebugPrinter(s.T(), false))

	for _, username := range s.testUsersName {
		registerResp := e.POST("/douyin/user/register/").
			WithQuery("username", username).WithQuery("password", s.password).
			WithFormField("username", username).WithFormField("password", s.password).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		registerResp.Value("status_code").Number().IsEqual(0)
		registerResp.Value("user_id").Number().Gt(0)
		registerResp.Value("token").String().Length().Gt(0)

		loginResp := e.POST("/douyin/user/login/").
			WithQuery("username", username).WithQuery("password", s.password).
			WithFormField("username", username).WithFormField("password", s.password).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginResp.Value("status_code").Number().IsEqual(0)
		loginResp.Value("user_id").Number().Gt(0)
		loginResp.Value("token").String().Length().Gt(0)

		token := loginResp.Value("token").String().Raw()
		userResp := e.GET("/douyin/user/").
			WithQuery("token", token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		userResp.Value("status_code").Number().IsEqual(0)
		userInfo := userResp.Value("user").Object()
		userInfo.NotEmpty()
		userInfo.Value("id").Number().Gt(0)
		userInfo.Value("name").String().Length().Gt(0)
		userInfo.Value("total_favorited").Number().Ge(0)
		userInfo.Value("work_count").Number().Ge(0)
		userInfo.Value("favorite_count").Number().Ge(0)
	}
}

func (s *BaseAPITestSuite) testPublish() {
	e := s.newHTTPExpect(httpexpect.NewDebugPrinter(s.T(), false))
	for i, username := range s.testUsersName {
		userId, token := getTestUserToken(username, s.password, e)

		for _, testVideoPath := range s.testVideoFilesPath[i] {
			// 标题为用户名+文件名，例如 "A1" "AZ14"
			title := username + filepath.Base(testVideoPath)
			publishResp := e.POST("/douyin/publish/action/").
				WithMultipart().
				WithFile("data", testVideoPath).
				WithFormField("token", token).
				WithFormField("title", title).
				Expect().
				Status(http.StatusOK).
				JSON().Object()
			publishResp.Value("status_code").Number().IsEqual(0)
		}

		publishListResp := e.GET("/douyin/publish/list/").
			WithQuery("user_id", userId).WithQuery("token", token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		publishListResp.Value("status_code").Number().IsEqual(0)
		publishListResp.Value("video_list").Array().Length().IsEqual(len(s.testVideoFilesPath[i]))

		for _, element := range publishListResp.Value("video_list").Array().Iter() {
			video := element.Object()
			video.ContainsKey("id")
			video.ContainsKey("author")
			video.Value("play_url").String().NotEmpty()
			video.Value("cover_url").String().NotEmpty()
			video.Value("title").String().NotEmpty()
		}
	}
}

func (s *BaseAPITestSuite) testFeed() {
	e := s.newHTTPExpect(httpexpect.NewDebugPrinter(s.T(), false))
	counter := 0
	var feedResp *httpexpect.Object
	var nextTime int64 = -1
	for nextTime != 0 {
		if nextTime == -1 {
			feedResp = e.GET("/douyin/feed/").
				Expect().
				Status(http.StatusOK).JSON().Object()
		} else {
			feedResp = e.GET("/douyin/feed/").
				WithQuery("latest_time", nextTime).
				Expect().
				Status(http.StatusOK).JSON().Object()
		}
		feedResp.Value("status_code").Number().IsEqual(0)
		var feedRespRaw model.FeedResponse
		feedResp.Decode(&feedRespRaw)
		nextTime = feedRespRaw.NextTime
		if nextTime == 0 {
			feedResp.Value("video_list").Array().Length().IsEqual(0)
		} else {
			feedResp.Value("next_time").Number().IsInt(64).Le(nextTime)
			feedResp.Value("video_list").Array().Length().Le(controller.FeedLimit)
			for _, element := range feedResp.Value("video_list").Array().Iter() {
				video := element.Object()
				video.ContainsKey("id")
				video.ContainsKey("author")
				video.Value("play_url").String().NotEmpty()
				video.Value("cover_url").String().NotEmpty()
				video.Value("title").String().NotEmpty()
				video.Value("favorite_count").Number().Ge(0)
				video.Value("comment_count").Number().Ge(0)
				video.ContainsKey("is_favorite")
				counter++
			}
		}
	}
	s.Require().Equal(s.totalVideoCount, counter)
}

func (s *BaseAPITestSuite) TestAPI() {
	if !s.Run("testUserAction", s.testUserAction) {
		return
	}
	if !s.Run("testPublish", s.testPublish) {
		return
	}
	if !s.Run("testFeed", s.testFeed) {
		return
	}
}

func TestBaseAPI(t *testing.T) {
	suite.Run(t, new(BaseAPITestSuite))
}
