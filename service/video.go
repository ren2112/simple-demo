package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"time"
)

func ToRespVideo(video model.Video) model.RespVideo {
	respVideo := model.RespVideo{
		Id:            video.Id,
		Author:        ToRespUser(video.Author),
		Title:         video.Title,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
	}
	return respVideo
}

func FeedVideoList(latestTimeUTC time.Time) (videoList []model.Video, err error) {
	err = common.DB.Preload("Author").
		Order("created_at DESC").
		Model(&model.Video{}).
		Where("created_at < ?", latestTimeUTC).
		Find(&videoList).Error
	return videoList, err
}

func JudgeFavorite(userId int64, videoId int64) (bool, error) {
	var isFavorite bool
	err := common.DB.Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Pluck("is_favorite", &isFavorite).Error
	return isFavorite, err
}

func JudgeRelation(authorId, userId int64) error {
	var follow model.Follow
	err := common.DB.Where("user_id = ? AND follower_user_id = ?", authorId, userId).First(&follow).Error
	return err
}

func PublishVideo(video model.Video, author model.User) error {
	var err error

	// 开始事务
	tx := common.DB.Begin()

	// 创建视频
	if err = tx.Create(&video).Error; err != nil {
		// 如果创建视频时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		return err
	}

	// 更新作者的work_count字段
	author.WorkCount++

	// 使用UpdateColumn更新作品计数字段
	if err = tx.Model(&author).UpdateColumn("work_count", author.WorkCount).Error; err != nil {
		// 如果更新作者信息时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		return err
	}

	// 提交事务
	tx.Commit()

	return err
}

func GetPublishVideoList(userId int64) ([]model.RespVideo, error) {
	var videoList []model.Video
	RespVideoList := []model.RespVideo{}
	err := common.DB.Preload("Author").Model(&videoList).Where("author_id=?", userId).Find(&videoList).Error
	//转化为响应专用
	for _, v := range videoList {
		RespVideoList = append(RespVideoList, ToRespVideo(v))
	}
	return RespVideoList, err
}
