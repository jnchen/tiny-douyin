package db

import (
	"douyin/model"
	"time"

	"gorm.io/gorm"
)

// Comment 评论 用户-视频
type Comment struct {
	ID        int64     `gorm:"primary_key;type:bigint;auto_increment;comment:评论id"`
	UserID    int64     `gorm:"type:bigint;comment:用户id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;association_autoupdate:false;association_autocreate:false;comment:用户信息"`
	VideoID   int64     `gorm:"type:bigint;comment:视频id"`
	Video     Video     `gorm:"foreignKey:VideoID;references:ID;association_autoupdate:false;association_autocreate:false;comment:视频信息"`
	Content   string    `gorm:"type:varchar(500);comment:评论内容"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt
}

// ToModel 转换为model.Comment，请确保User不为空。
func (comment *Comment) ToModel() *model.Comment {
	return &model.Comment{
		Id:         comment.ID,
		User:       *comment.User.ToModel(),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("01-02"),
	}
}

// AfterCreate 插入新的评论信息后，需要更新视频的评论数
func (comment *Comment) AfterCreate(tx *gorm.DB) (err error) {
	// 更新视频的评论数
	result := tx.Model(&Video{}).
		Where("id = ?", comment.VideoID).
		Update(
			"comment_count",
			gorm.Expr("comment_count + ?", 1),
		)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return
}

// BeforeDelete 删除评论信息前，需要更新视频的评论数
func (comment *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	// 更新视频的评论数
	result := tx.Model(&Video{}).
		Where("id = ?", comment.VideoID).
		Update(
			"comment_count",
			gorm.Expr("comment_count - ?", 1),
		)
	if nil != result.Error {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return
}
