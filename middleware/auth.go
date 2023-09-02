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
		token, exists := c.GetQuery("token")
		if !exists {
			token, exists = c.GetPostForm("token")
		}
		token = strings.TrimSpace(token)

		// 如果不需要token，直接跳过
		if !isTokenRequired && token == "" {
			c.Next()
			return
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  "请登录！",
			})
			return
		}

		user, err := service.CheckLogin(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
