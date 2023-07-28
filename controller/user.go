package controller

import (
	"douyin/model"
	"douyin/service"
	"douyin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

// Register 用户注册
// c: http上下文
// return: 用户id和用户token
func Register(c *gin.Context) {
	// 用户名
	username := c.Query("username")
	// 密码
	password := c.Query("password")

	token := util.UUIDNoLine()

	exist, err := service.UserExists(username)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	if exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}

	passwordMd5, err := util.Md5(password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
		return
	}

	newUser, err := service.UserCreate(username, passwordMd5)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "注册失败"},
		})
		return
	}

	timeDuration, err := time.ParseDuration("4h")
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "生成Token失败"},
		})
		return
	}

	expireAt := time.Now().Add(timeDuration)

	err = service.UserTokenCreate(newUser.ID, token, expireAt)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "生成Token失败"},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   newUser.ID,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := service.CheckLogin(token); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	user, exist := service.CheckLogin(token)
	fmt.Println(user)
	if exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     model.User{},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
