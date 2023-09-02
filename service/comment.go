package service

import (
	"douyin/db"
	"errors"
)

func CommentPost(userId int64, videoId int64, content string) (*db.Comment, error) {
	var comment = db.Comment{
		UserID:  userId,
		VideoID: videoId,
		Content: content,
	}
	// 创建后加载对象，方便直接调用ToModel方法
	result := db.ORM().Preload("User").Create(&comment).First(&comment)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("sql execution failed")
	}
	return &comment, nil
}

func CommentDelete(commentId, videoID int64) error {
	var comment = db.Comment{
		ID:      commentId,
		VideoID: videoID,
	}
	result := db.ORM().Delete(&comment)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("sql execution failed")
	}
	return nil
}

func CommentList(videoId int64) ([]db.Comment, error) {
	var commentList []db.Comment
	result := db.ORM().Preload("User").Where("video_id = ?", videoId).Find(&commentList)
	if nil != result.Error {
		return nil, result.Error
	}
	return commentList, nil
}
