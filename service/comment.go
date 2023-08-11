package service

import (
	"douyin/db"
	"errors"
)

func CommentPost(userId int64, videoId int64, content string) (*db.Comment, error) {
	var comment db.Comment = db.Comment{
		UserID:  userId,
		VideoID: videoId,
		Content: content,
	}
	result := db.DB.Create(&comment)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("sql execution failed")
	}
	return &comment, nil
}

func CommentDelete(commentId int64) error {
	var comment db.Comment = db.Comment{
		ID: commentId,
	}
	result := db.DB.Delete(&comment)
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
	result := db.DB.Preload("User").Where("video_id = ?", videoId).Find(&commentList)
	if nil != result.Error {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("sql execution failed")
	}
	return commentList, nil
}
