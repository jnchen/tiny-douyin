package db

import (
	"time"
)

// Follow 关注 用户-用户
type Follow struct {
	ID         int64     `gorm:"primary_key;type:bigint;auto_increment;comment:关注id"`
	FollowerID int64     `gorm:"type:bigint;comment:关注者id"`
	FollowedID int64     `gorm:"type:bigint;comment:被关注者id"`
	CreatedAt  time.Time `gorm:"comment:创建时间"`
	UpdatedAt  time.Time `gorm:"comment:更新时间"`
}
