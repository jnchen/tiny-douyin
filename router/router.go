package router

import (
	"douyin/controller"
	"douyin/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(r *gin.Engine) error {
	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return fmt.Errorf("设置信任代理失败：%w", err)
	}

	r.Use(gin.Logger(), gin.Recovery())

	// public directory is used to serve static resources
	r.Static("/static", "./public")
	r.LoadHTMLGlob("templates/*")

	// home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET(
		"/user/",
		middleware.Auth(true),
		controller.UserInfo,
	)
	apiRouter.POST(
		"/user/register/",
		// middleware.QPSLimit(3),
		controller.Register,
	)
	apiRouter.POST(
		"/user/login/",
		// middleware.QPSLimit(3),
		controller.Login,
	)
	apiRouter.POST(
		"/publish/action/",
		middleware.Auth(true),
		controller.Publish,
	)
	apiRouter.GET(
		"/publish/list/",
		middleware.Auth(false),
		controller.PublishList,
	)

	// extra apis - I
	apiRouter.POST(
		"/favorite/action/",
		middleware.Auth(true),
		controller.FavoriteAction,
	)
	apiRouter.GET(
		"/favorite/list/",
		middleware.Auth(false),
		controller.FavoriteList,
	)
	apiRouter.POST(
		"/comment/action/",
		// middleware.QPSLimit(1000),
		middleware.Auth(true),
		controller.CommentAction,
	)
	apiRouter.GET(
		"/comment/list/",
		middleware.Auth(false),
		controller.CommentList,
	)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)

	return nil
}
