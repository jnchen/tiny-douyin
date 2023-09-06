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
	result := db.ORM().
		Preload("User").
		Create(&comment).
		First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("发表评论失败")
	}
	return &comment, nil
}

func CommentDelete(commentId, videoID int64) (err error) {
	var comment = db.Comment{
		ID:      commentId,
		VideoID: videoID,
	}
	tx := db.ORM().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result := tx.Delete(&comment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除评论失败")
	}
	return nil
}

func CommentList(videoId int64) ([]db.Comment, error) {
	var commentList []db.Comment
	result := db.ORM().
		Preload("User").
		Where("video_id = ?", videoId).
		Order("created_at DESC").
		Find(&commentList)
	if nil != result.Error {
		return nil, result.Error
	}
	return commentList, nil
}
