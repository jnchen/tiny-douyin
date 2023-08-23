package router

import (
	"douyin/controller"
	"douyin/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine) {
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
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
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
}
