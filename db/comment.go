package db

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论 用户-视频
type Comment struct {
	ID        int64     `gorm:"primary_key;type:bigint;auto_increment;comment:评论id"`
	UserID    int64     `gorm:"type:bigint;comment:用户id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;association_autoupdate:false;association_autocreate:false;comment:用户信息"`
	VideoID   int64     `gorm:"type:bigint;comment:视频id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID;association_autoupdate:false;association_autocreate:false;comment:视频信息"`
	Content   string    `gorm:"type:varchar(500);comment:评论内容"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt
}
