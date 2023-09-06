package db

import (
	"douyin/model"
	"time"
)

// User 用户
type User struct {
	ID              int64     `gorm:"primary_key;type:bigint;auto_increment;comment:用户id"`
	Username        string    `gorm:"type:varchar(32);uniqueIndex;comment:用户登录名"`
	Password        string    `gorm:"type:varchar(64);comment:用户密码"`
	Avatar          string    `gorm:"type:varchar(512);comment:头像地址"`
	BackgroundImage string    `gorm:"type:varchar(1024);comment:背景地址"`
	Signature       string    `gorm:"type:varchar(1024);comment:个性签名"`
	TotalFavorited  int64     `gorm:"type:bigint;comment:获赞总数"`
	WorkCount       int64     `gorm:"type:bigint;comment:作品数"`
	FavoriteCount   int64     `gorm:"type:bigint;comment:喜欢数"`
	CreatedAt       time.Time `gorm:"comment:创建时间"`
	UpdatedAt       time.Time `gorm:"comment:更新时间"`
}

func (user *User) ToModel() *model.User {
	return &model.User{
		Id:              user.ID,
		Name:            user.Username,
		FollowCount:     0, // TODO
		FollowerCount:   0, // TODO
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}
