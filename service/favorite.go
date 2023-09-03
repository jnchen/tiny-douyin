package service

import (
	"douyin/db"
	"errors"
)

func FavoriteAction(userId int64, videoId int64) error {
	favorite := db.Favorite{
		UserID:  userId,
		VideoID: videoId,
	}
	result := db.ORM().Create(&favorite)
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
	result := db.ORM().
		Delete(&favorite)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("取消点赞失败")
	}
	return nil
}

func FavoriteList(userId int64) ([]db.Video, error) {
	var videoList []db.Video
	result := db.ORM().
		Preload("Author").
		Joins("JOIN favorite ON favorite.video_id = video.id").
		Where("favorite.user_id = ?", userId).
		Find(&videoList)
	if nil != result.Error {
		return nil, result.Error
	}
	return videoList, nil
}
