package test

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type InteractAPITestSuite struct {
	APITestSuite
}

func (s *InteractAPITestSuite) testFavorite() {
	e := s.newHTTPExpect(httpexpect.NewDebugPrinter(s.T(), false))

	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().Gt(0)
	firstVideo := feedResp.Value("video_list").Array().Value(0).Object()
	videoId := firstVideo.Value("id").Number().Raw()

	userId, token := getTestUserToken(s.testUserA, s.password, e)

	favoriteResp := e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	favoriteResp.Value("status_code").Number().IsEqual(0)
	favoriteResp = e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	favoriteResp.Value("status_code").Number().IsEqual(1)
	favoriteResp.Value("status_msg").String().IsEqual("重复点赞")

	favoriteListResp := e.GET("/douyin/favorite/list/").
		WithQuery("token", token).WithQuery("user_id", userId).
		WithFormField("token", token).WithFormField("user_id", userId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	favoriteListResp.Value("status_code").Number().IsEqual(0)
	for _, element := range favoriteListResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}

	favoriteUndoResp := e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 2).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	favoriteUndoResp.Value("status_code").Number().IsEqual(0)
	favoriteUndoResp = e.POST("/douyin/favorite/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 2).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	favoriteUndoResp.Value("status_code").Number().IsEqual(1)
	favoriteUndoResp.Value("status_msg").String().IsEqual("取消点赞失败")
}

func (s *InteractAPITestSuite) testComment() {
	t := s.T()
	e := s.newHTTPExpect(httpexpect.NewDebugPrinter(t, false))

	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().IsEqual(0)
	feedResp.Value("video_list").Array().Length().Gt(0)
	firstVideo := feedResp.Value("video_list").Array().Value(0).Object()
	videoId := firstVideo.Value("id").Number().Raw()

	_, token := getTestUserToken(s.testUserB, s.password, e)

	addCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).WithQuery("comment_text", "测试评论").
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 1).WithFormField("comment_text", "测试评论").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	addCommentResp.Value("status_code").Number().IsEqual(0)
	addCommentResp.Value("comment").Object().Value("id").Number().Gt(0)
	commentId := int(addCommentResp.Value("comment").Object().Value("id").Number().Raw())

	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		WithFormField("token", token).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)
	containTestComment := false
	for _, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		comment.ContainsKey("id")
		comment.ContainsKey("user")
		comment.Value("content").String().NotEmpty()
		comment.Value("create_date").String().NotEmpty()
		if int(comment.Value("id").Number().Raw()) == commentId {
			containTestComment = true
		}
	}

	assert.True(t, containTestComment, "Can't find test comment in list")

	delCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", commentId).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 2).WithFormField("comment_id", commentId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	delCommentResp.Value("status_code").Number().IsEqual(0)
	delCommentResp = e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		WithQuery("action_type", 2).WithQuery("comment_id", commentId).
		WithFormField("token", token).WithFormField("video_id", videoId).
		WithFormField("action_type", 2).WithFormField("comment_id", commentId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	delCommentResp.Value("status_code").Number().IsEqual(1)
	delCommentResp.Value("status_msg").String().IsEqual("删除评论失败")
}

func (s *InteractAPITestSuite) TestAPI() {
	if !s.Run("testFavorite", s.testFavorite) {
		return
	}
	if !s.Run("testComment", s.testComment) {
		return
	}
}

func TestInteractAPI(t *testing.T) {
	suite.Run(t, new(InteractAPITestSuite))
}
