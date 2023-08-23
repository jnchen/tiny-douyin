package util

import (
	"douyin/model"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) *model.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	} else {
		return user.(*model.User)
	}
}
