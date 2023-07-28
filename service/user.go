package service

import (
	"douyin/db"
	"douyin/util"
	"errors"
	"time"
)

func UserExists(username string) (bool, error) {
	var count int64
	result := db.DB.Model(&db.User{}).Where("username = ?", username).Count(&count)
	if nil != result.Error {
		return false, result.Error
	}
	if 0 == result.RowsAffected {
		return false, nil
	}
	return count > 0, nil
}

func UserCreate(username string, password string) (*db.User, error) {
	passwordMd5, err := util.Md5(password)
	if err != nil {
		return nil, err
	}
	user := db.User{
		Name:      username,
		Username:  username,
		Password:  passwordMd5,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	result := db.DB.Create(&user)
	if nil != result.Error {
		return nil, result.Error
	}
	if 0 == result.RowsAffected {
		return nil, errors.New("sql执行失败")
	}
	return &user, nil
}
