package service

import (
	"douyin/db"
	"gorm.io/gorm"
)

func FavoriteCheck(userId int64, videoId int64) (bool, error) {
	var favorite db.Favorite
	result := db.DB.
		Where("user_id = ? and video_id = ?", userId, videoId).
		Limit(1).
		Find(&favorite)
	if nil != result.Error {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func FavoriteAction(userId int64, videoId int64) error {
	favorite := db.Favorite{
		UserID:  userId,
		VideoID: videoId,
	}
	result := db.DB.Create(&favorite)
	if nil != result.Error {
		return result.Error
	}
	return nil
}

func FavoriteDelete(userId int64, videoId int64) error {
	favorite := db.Favorite{
		UserID:  userId,
		VideoID: videoId,
	}
	result := db.DB.
		Where("user_id = ? and video_id = ?", favorite.UserID, favorite.VideoID).
		Delete(&favorite)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func FavoriteList(userId int64) ([]db.Video, error) {
	var videoList []db.Video
	result := db.DB.Preload("Author").
		Joins("JOIN favorite ON favorite.video_id = video.id").
		Where("favorite.user_id = ?", userId).
		Find(&videoList)
	if nil != result.Error {
		return nil, result.Error
	}
	return videoList, nil
}
