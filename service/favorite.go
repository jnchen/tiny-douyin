package service

import (
	"douyin/db"
	"errors"
	"github.com/go-sql-driver/mysql"
)

func FavoriteAction(userId int64, videoId int64) error {
	favorite := db.Favorite{
		UserID:  userId,
		VideoID: videoId,
	}
	result := db.ORM().Create(&favorite)
	var mysqlErr *mysql.MySQLError
	errors.As(result.Error, &mysqlErr)
	if mysqlErr != nil {
		if mysqlErr.Number == 1062 {
			return errors.New("重复点赞")
		}
		return result.Error
	}
	return nil
}

func FavoriteDelete(userId int64, videoId int64) (err error) {
	favorite := db.Favorite{
		UserID:  userId,
		VideoID: videoId,
	}
	tx := db.ORM().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result := tx.Delete(&favorite)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 { // 未找到记录，回滚 `(*Favorite).BeforeDelete` 的更新
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
