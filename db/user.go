package db

import (
	"time"
)

// User 用户
type User struct {
	ID              int64     `gorm:"primary_key;type:bigint;auto_increment;comment:用户id"`
	Name            string    `gorm:"type:varchar(200);comment:用户名称"`
	Avatar          string    `gorm:"type:varchar(500);comment:头像地址"`
	BackgroundImage string    `gorm:"type:varchar(500);comment:背景地址"`
	Signature       string    `gorm:"type:varchar(1000);comment:个性签名"`
	CreatedAt       time.Time `gorm:"comment:创建时间"`
	UpdatedAt       time.Time `gorm:"comment:更新时间"`
}
