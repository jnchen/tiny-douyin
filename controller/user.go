package controller

import (
	"douyin/model"
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 用户注册
// c: http上下文
// return: 用户id和用户token
func Register(c *gin.Context) {
	// 用户名
	username := c.Query("username")
	// 密码
	password := c.Query("password")

	exist, err := service.UserExists(username)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	if exist {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}

	newUser, err := service.UserCreate(username, password)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "注册失败"},
		})
		return
	}

	token, err := service.UserTokenCreate(newUser.ID)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "生成Token失败"},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   newUser.ID,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userId, err := service.UserLogin(username, password)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Login Error"},
		})
		return
	}

	token, err := service.UserTokenCreate(userId)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Login Error"},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	user, exist := service.CheckLogin(token)
	fmt.Println(user)
	if exist {
		c.JSON(http.StatusOK, model.UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     *user,
		})
	} else {
		c.JSON(http.StatusOK, model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
