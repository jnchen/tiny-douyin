package db

import (
	"time"
)

// Video 视频
type Video struct {
	ID        int64     `gorm:"primary_key;auto_increment;type:bigint;comment:视频id"`
	AuthorId  int64     `gorm:"type:bigint;comment:作者id"`
	Author    User      `gorm:"foreignKey:AuthorId;references:ID;association_autoupdate:false;association_autocreate:false;comment:投稿者信息"`
	PlayUrl   string    `gorm:"type:varchar(2048);comment:播放地址"`
	CoverUrl  string    `gorm:"type:varchar(2048);comment:封面图像地址"`
	Title     string    `gorm:"type:varchar(128);comment:视频标题"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
