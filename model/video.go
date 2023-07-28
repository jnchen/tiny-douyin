package model

import (
	"time"
)

// Video 视频
type Video struct {
	ID        int64     `gorm:"primary_key;auto_increment;type:bigint;comment:视频id"`
	AuthorId  int64     `gorm:"type:bigint;comment:作者id"`
	PlayUrl   string    `gorm:"type:varchar(500);comment:播放地址"`
	CoverUrl  string    `gorm:"type:varchar(500);comment:封面图像地址"`
	Title     string    `gorm:"type:varchar(100);comment:视频标题"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
