package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"testing"
)

func TestFeed(t *testing.T) {
	e := newExpect(t)

	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().Gt(0)

	for _, element := range feedResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}

func TestUserAction(t *testing.T) {
	e := newExpect(t)

	// rand.Seed(time.Now().UnixNano())
	registerValue := fmt.Sprintf("douyin%d", rand.Intn(65536))

	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	registerResp.Value("status_code").Number().IsEqual(0)
	registerResp.Value("user_id").Number().Gt(0)
	registerResp.Value("token").String().Length().Gt(0)

	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
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
}

func TestPublish(t *testing.T) {
	e := newExpect(t)

	testVideoDir := filepath.Join("../public", "test_videos")
	testUserPaths, err := filepath.Glob(filepath.Join(testVideoDir, "[A-Z]*"))
	if nil != err {
		t.Fatal(err)
	}
	for _, testUserPath := range testUserPaths {
		letters := filepath.Base(testUserPath)
		userId, token := getTestUserToken(letters, e)

		testVideoPaths, err := filepath.Glob(filepath.Join(
			testVideoDir,
			letters,
			"[0-9]*.mp4",
		))
		if nil != err {
			t.Fatal(err)
		}

		t.Log("Testing publish")
		for _, testVideoPath := range testVideoPaths {
			title := letters + filepath.Base(testVideoPath)
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
		t.Log("Tested publish")

		t.Log("Testing publish list")
		publishListResp := e.GET("/douyin/publish/list/").
			WithQuery("user_id", userId).WithQuery("token", token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		publishListResp.Value("status_code").Number().IsEqual(0)
		publishListResp.Value("video_list").Array().Length().IsEqual(len(testVideoPaths))

		for _, element := range publishListResp.Value("video_list").Array().Iter() {
			video := element.Object()
			video.ContainsKey("id")
			video.ContainsKey("author")
			video.Value("play_url").String().NotEmpty()
			video.Value("cover_url").String().NotEmpty()
			video.Value("title").String().NotEmpty()
		}
		t.Log("Tested publish list")
	}
}
