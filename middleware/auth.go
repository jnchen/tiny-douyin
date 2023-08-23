package middleware

import (
	"douyin/model"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(isTokenRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth model.Auth
		if err := c.ShouldBind(&auth); err != nil {
			c.JSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			c.Abort()
			return
		}
		token := strings.TrimSpace(auth.Token)

		// 如果不需要token，直接跳过
		if !isTokenRequired && token == "" {
			c.Next()
			return
		}

		if token == "" {
			c.JSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  "请登录！",
			})
			c.Abort()
			return
		}

		user, err := service.CheckLogin(token)
		if err != nil {
			c.JSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
