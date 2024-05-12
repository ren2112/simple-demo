package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func FavoriteAction(actionType string, userId int64, videoId int64) error {
	var favorite model.Favorite
	favorite.UserId = userId
	favorite.VideoId = videoId

	//开启事务
	tx := common.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var author model.User
	var video model.Video

	//获得视频作者的信息，因为要对视频作者的获赞数量更改
	common.DB.Preload("Author").Where("id=?", videoId).First(&video)
	author = video.Author
	if actionType == "1" {
		favorite.IsFavorite = true
		// 查看favorite表格是否存在数据，如果已经存在且已经点过赞，则直接返回
		var existingFavorite model.Favorite
		if result := tx.Where("user_id = ? AND video_id = ?", userId, videoId).First(&existingFavorite); result.RowsAffected == 0 {
			tx.Create(&favorite)
		} else if existingFavorite.IsFavorite == true {
			return errors.New("请勿重复点赞！")
		} else { //将false更新为true
			if err := tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Update("is_favorite", true).Error; err != nil {
				return err
			}
		}

		// 更新视频的favorite_count字段
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新视频作者的TotalFavorited字段
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新用户的favorite_count字段
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}
	} else {
		//查找是否存在点赞记录，若没有或者已经是没点赞的状态下则无法取消点赞
		var existFavorite model.Favorite
		tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).First(&existFavorite)
		if existFavorite.Id == 0 || existFavorite.IsFavorite == false {
			return errors.New("请勿在没点赞情况下取消点赞")
		}

		tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Update("is_favorite", false)
		// 更新视频的favorite_count字段
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新视频作者的TotalFavorited字段
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新用户的favorite_count字段
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func GetFavoriteList(userId string) ([]model.RespVideo, error) {
	var favorites []model.Favorite
	var resVideoList []model.RespVideo
	if err := common.DB.Preload("Video").Preload("Video.Author").Model(&model.Favorite{}).Where("user_id = ?", userId).Find(&favorites).Error; err != nil {
		return nil, err
	}

	//设置视频结构体为为喜欢并且转换为响应结构体
	for _, v := range favorites {
		//v.Video.IsFavorite = true
		if v.IsFavorite {
			resVideoList = append(resVideoList, ToRespVideo(v.Video))
		}
	}
	return resVideoList, nil
}
