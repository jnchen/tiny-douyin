package db

import (
	"gorm.io/gorm"
	"time"
)

// UserToken 用户token表，正常应该用redis什么的，这里直接用数据库存
type UserToken struct {
	gorm.Model
	UserId   int64     `gorm:"type:bigint;comment:用户id"`
	Token    string    `gorm:"type:varchar(200);comment:用户token"`
	ExpireAt time.Time `gorm:"type:datetime;comment:过期时间"`
}
