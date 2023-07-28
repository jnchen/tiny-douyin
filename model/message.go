package model

import (
	"time"
)

// Message 消息 用户-用户
type Message struct {
	ID         int64     `gorm:"primary_key;type:bigint;auto_increment;comment:消息id"`
	ToUserID   int64     `gorm:"type:bigint;comment:消息接收者id"`
	FromUserID int64     `gorm:"type:bigint;comment:消息发送者id"`
	Content    string    `gorm:"type:varchar(500);comment:消息内容"`
	CreatedAt  time.Time `gorm:"comment:创建时间"`
	UpdatedAt  time.Time `gorm:"comment:更新时间"`
}
