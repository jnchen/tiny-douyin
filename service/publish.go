package service

import (
	"douyin/db"
	"errors"
	"fmt"
	"time"
)

func VideoPublish(
	userId int64,
	videoUrl, coverUrl,
	title string,
) (*db.Video, error) {
	video := db.Video{
		AuthorId: userId,
		PlayUrl:  videoUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}
	// 创建后加载对象，方便直接调用ToModel方法
	result := db.DB.Preload("Author").Create(&video).First(&video)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("发布视频失败！")
	}
	return &video, nil
}

func VideoPublishList(userId int64) ([]db.Video, error) {
	var videoPublishList []db.Video
	result := db.DB.Preload("Author").
		Where("author_id = ?", userId).
		Find(&videoPublishList)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("获取视频发布列表失败！")
	}
	return videoPublishList, nil
}

func VideoList(userId int64, latestTime time.Time, limit int) ([]db.VideoWithFavorite, error) {
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
	fmt.Println(userId)
	for i, video := range videoList {
		fmt.Println(i, video.IsFavorite)
	}
	return videoList, nil
}
