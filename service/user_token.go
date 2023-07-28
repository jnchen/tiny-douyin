package service

import (
	"douyin/db"
	"douyin/model"
	"douyin/util"
	"time"
)

func UserTokenCreate(id int64) (string, error) {
	token := util.UUIDNoLine()
	timeDuration, err := time.ParseDuration("4h")
	if err != nil {
		return "", err
	}
	expireAt := time.Now().Add(timeDuration)

	db.DB.Create(&db.UserToken{
		UserId:   id,
		Token:    token,
		ExpireAt: expireAt,
	})

	return token, nil
}

func CheckLogin(token string) (*model.User, bool) {
	var userToken db.UserToken
	result := db.DB.Where("token = ?", token).First(&userToken)
	if result.Error != nil {
		return nil, false
	}
	if 0 == result.RowsAffected {
		return nil, false
	}

	var userInfo db.User
	result = db.DB.Where("id= ?", userToken.UserId).First(&userInfo)
	if result.Error != nil {
		return nil, false
	}
	if 0 == result.RowsAffected {
		return nil, false
	}

	return &model.User{
		Id:            userInfo.ID,
		Name:          userInfo.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}, true
}

func UserLogin(username string, password string) (int64, error) {
	return 0, nil
}
