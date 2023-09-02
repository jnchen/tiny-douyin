package db

import (
	"gorm.io/gorm"
	"time"
)

// Favorite 点赞 用户-视频
type Favorite struct {
	// ID        int64     `gorm:"primary_key;type:bigint;auto_increment;comment:点赞id"`
	UserID    int64     `gorm:"primary_key;type:bigint;comment:用户id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;association_autoupdate:false;association_autocreate:false;comment:用户信息"`
	VideoID   int64     `gorm:"primary_key;type:bigint;comment:视频id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID;association_autoupdate:false;association_autocreate:false;comment:视频信息"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
}

// AfterCreate 插入新的点赞信息后，需要更新视频的点赞数、用户的点赞数、用户的被点赞数
func (favorite *Favorite) AfterCreate(tx *gorm.DB) (err error) {
	// 更新视频的点赞数
	if err = tx.Model(&Video{}).
		Where("id = ?", favorite.VideoID).
		Update("favorite_count", gorm.Expr(
			"favorite_count + ?",
			1,
		)).Error; err != nil {
		return
	}

	// 更新用户的点赞数
	if err = tx.Model(&User{}).
		Where("id = ?", favorite.UserID).
		Update("favorite_count", gorm.Expr(
			"favorite_count + ?",
			1,
		)).Error; err != nil {
		return
	}

	// 更新用户的被点赞数
	favorite.Video.ID = favorite.VideoID
	if err = orm.Select("author_id").
		First(&favorite.Video).Error; err != nil {
		return
	}
	if err = tx.Model(&User{}).
		Where("id = ?", favorite.Video.AuthorId).
		Update("total_favorited", gorm.Expr(
			"total_favorited + ?",
			1,
		)).Error; err != nil {
		return
	}

	return
}

// BeforeDelete 删除点赞信息前，需要更新视频的点赞数、用户的点赞数、用户的被点赞数
func (favorite *Favorite) BeforeDelete(tx *gorm.DB) (err error) {
	// 更新视频的点赞数
	if err = tx.Model(&Video{}).
		Where("id = ?", favorite.VideoID).
		Update("favorite_count", gorm.Expr(
			"favorite_count - ?",
			1,
		)).Error; err != nil {
		return
	}

	// 更新用户的点赞数
	if err = tx.Model(&User{}).
		Where("id = ?", favorite.UserID).
		Update("favorite_count", gorm.Expr(
			"favorite_count - ?",
			1,
		)).Error; err != nil {
		return
	}

	// 更新用户的被点赞数
	favorite.Video.ID = favorite.VideoID
	if err = orm.Select("author_id").
		First(&favorite.Video).Error; err != nil {
		return
	}
	if err = tx.Model(&User{}).
		Where("id = ?", favorite.Video.AuthorId).
		Update("total_favorited", gorm.Expr(
			"total_favorited - ?",
			1,
		)).Error; err != nil {
		return
	}
	return
}
