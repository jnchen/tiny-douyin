package service

import (
	"douyin/db"
	"errors"
)

func PublishAction(
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
	result := db.ORM().
		Preload("Author").
		Create(&video).
		Limit(1).
		Find(&video)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("发布视频失败！")
	}
	return &video, nil
}

func PublishDelete(videoId int64) error {
	result := db.ORM().Delete(&db.Video{}, videoId)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除视频失败！")
	}
	return nil
}

func PublishList(userId int64) ([]db.Video, error) {
	var videoPublishList []db.Video
	result := db.ORM().Preload("Author").
		Where("author_id = ?", userId).
		Find(&videoPublishList)
	if nil != result.Error {
		return nil, result.Error
	}
	return videoPublishList, nil
}
