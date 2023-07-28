package model

import (
	"time"
)

// Favorite 点赞 用户-视频
type Favorite struct {
	ID        int64     `gorm:"primary_key;type:bigint;auto_increment;comment:点赞id"`
	UserID    int64     `gorm:"type:bigint;comment:用户id"`
	VideoID   int64     `gorm:"type:bigint;comment:视频id"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
