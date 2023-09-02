package service

import (
	"douyin/db"
	"douyin/util"
	"errors"
	"fmt"
	"strings"
	"time"
)

func UserExists(username string) (bool, error) {
	result := db.ORM().Where("username = ?", username).Find(&db.User{})
	if nil != result.Error {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func UserCreate(username string, password string) (*db.User, error) {
	usernameMd5, err := util.Md5(strings.ToLower(username))
	if err != nil {
		return nil, err
	}
	passwordMd5, err := util.Md5(password)
	if err != nil {
		return nil, err
	}
	user := db.User{
		Username: username,
		Password: passwordMd5,
		Avatar: fmt.Sprintf(
			"https://avatar.marktion.cn/api/avatar/%s?t=github",
			usernameMd5,
		),
		BackgroundImage: util.RandomImageURL([]string{}),
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}
	result := db.ORM().Create(&user)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("sql执行失败")
	}
	return &user, nil
}
