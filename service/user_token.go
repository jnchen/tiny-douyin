package service

import (
	"douyin/model"
	"time"
)

func UserTokenCreate(id int64, token string, expireAt time.Time) error {
	return nil
}

func CheckLogin(token string) (*model.User, bool) {
	return nil, true
}
