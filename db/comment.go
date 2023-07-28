package db

import (
	"time"
)

// Comment 评论 用户-视频
type Comment struct {
	ID        int64     `gorm:"primary_key;type:bigint;auto_increment;comment:评论id"`
	UserID    int64     `gorm:"type:bigint;comment:用户id"`
	VideoID   int64     `gorm:"type:bigint;comment:视频id"`
	Content   string    `gorm:"type:varchar(500);comment:评论内容"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}
