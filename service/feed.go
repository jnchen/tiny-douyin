package service

import (
	"douyin/db"
	"time"
)

func FeedList(userId int64, latestTime time.Time, limit int) ([]db.VideoWithFavorite, error) {
	var videoList []db.VideoWithFavorite
	result := db.DB.
		Model(&db.Video{}).
		Preload("Author").
		Select("video.*, favorite.user_id IS NOT NULL AS is_favorite").
		Joins("LEFT JOIN favorite ON video.id = favorite.video_id AND favorite.user_id = ?", userId).
		Where("video.created_at < ?", latestTime).
		Order("video.created_at DESC").
		Limit(limit).
		Find(&videoList)
	if nil != result.Error {
		return nil, result.Error
	}
	return videoList, nil
}
