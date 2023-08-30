package db

import (
	"douyin/model"
	"gorm.io/gorm"
	"time"
)

// Video 视频
type Video struct {
	ID            int64     `gorm:"primary_key;auto_increment;type:bigint;comment:视频id"`
	AuthorId      int64     `gorm:"type:bigint;comment:作者id"`
	Author        User      `gorm:"foreignKey:AuthorId;references:ID;association_autoupdate:false;association_autocreate:false;comment:投稿者信息"`
	PlayUrl       string    `gorm:"type:varchar(2048);comment:播放地址"`
	CoverUrl      string    `gorm:"type:varchar(2048);comment:封面图像地址"`
	FavoriteCount int64     `gorm:"type:bigint;comment:点赞数"`
	CommentCount  int64     `gorm:"type:bigint;comment:评论数"`
	Title         string    `gorm:"type:varchar(128);comment:视频标题"`
	CreatedAt     time.Time `gorm:"comment:创建时间"`
	UpdatedAt     time.Time `gorm:"comment:更新时间"`
}

type VideoWithFavorite struct {
	Video
	IsFavorite bool `gorm:"->;column:is_favorite;comment:是否已经收藏"`
}

// ToModel 转换为model.Video，请确保Author不为空。
func (video *Video) ToModel() *model.Video {
	return &model.Video{
		Id:            video.ID,
		Author:        *video.Author.ToModel(),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    false,
		Title:         video.Title,
	}
}

// ToModel 转换为model.Video，请确保Author不为空。
func (video *VideoWithFavorite) ToModel() *model.Video {
	return &model.Video{
		Id:            video.ID,
		Author:        *video.Author.ToModel(),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
		Title:         video.Title,
	}
}

// AfterCreate 插入新的视频信息后，需要更新作者的作品数
func (video *Video) AfterCreate(tx *gorm.DB) (err error) {
	// 更新作者的作品数
	err = tx.Model(&User{}).
		Where("id = ?", video.AuthorId).
		Update(
			"work_count",
			gorm.Expr("work_count + ?", 1),
		).Error
	return
}

// AfterDelete 删除视频信息后，需要更新作者的作品数
func (video *Video) AfterDelete(tx *gorm.DB) (err error) {
	// 更新作者的作品数
	err = tx.Model(&User{}).
		Where("id = ?", video.AuthorId).
		Update(
			"work_count",
			gorm.Expr("work_count - ?", 1),
		).Error
	return
}
