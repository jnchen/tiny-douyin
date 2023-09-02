package service

import (
	"douyin/db"
	"douyin/model"
	"douyin/util"
	"errors"
	"time"

	"gorm.io/gorm"
)

func UserTokenCreate(id int64) (string, error) {
	token := util.UUIDNoLine()
	timeDuration, err := time.ParseDuration("4h")
	if err != nil {
		return "", err
	}
	expireAt := time.Now().Add(timeDuration)

	if result := db.ORM().Create(&db.UserToken{
		UserId:   id,
		Token:    token,
		ExpireAt: expireAt,
	}); result.Error != nil {
		return "", result.Error
	}

	return token, nil
}

func CheckLogin(token string) (*model.User, error) {
	var userToken db.UserToken
	result := db.ORM().Where("token = ?", token).Limit(1).Find(&userToken)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("token不存在")
	}

	if time.Now().After(userToken.ExpireAt) {
		delCond := db.UserToken{
			Model: gorm.Model{
				ID: userToken.ID,
			},
		}
		db.ORM().Delete(&delCond)
		return nil, errors.New("token已过期")
	}

	var userInfo db.User
	result = db.ORM().Where("id = ?", userToken.UserId).First(&userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}

	return userInfo.ToModel(), nil
}

func UserLogin(username string, password string) (int64, error) {
	passwordMd5, err := util.Md5(password)
	if err != nil {
		return 0, err
	}
	var user db.User
	result := db.ORM().Where("username = ? and password = ?", username, passwordMd5).First(&user)
	if nil != result.Error {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("用户名或密码错误")
		}
		return 0, result.Error
	}

	return user.ID, nil
}
