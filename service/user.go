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
	var count int64
	result := db.DB.Model(&db.User{}).Where("username = ?", username).Count(&count)
	if nil != result.Error {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return count > 0, nil
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
		Name:     username,
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
	result := db.DB.Create(&user)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("sql执行失败")
	}
	return &user, nil
}
