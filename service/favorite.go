package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	redislock "github.com/jefferyjob/go-redislock"
	"gorm.io/gorm"
)

func FavoriteAction(actionType int32, userId int64, videoId int64) error {
	ctx := context.Background()
	var favorite model.Favorite
	favorite.UserId = userId
	favorite.VideoId = videoId

	//上redis分布式锁锁视频保证并发情况点赞安全
	lock := redislock.New(ctx, common.RedisClient, fmt.Sprintf("favorite_video:%d", videoId), redislock.WithAutoRenew())
	err := lock.Lock()
	if err != nil {
		return errors.New("操作失败！")
	}

	defer lock.UnLock()
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
	tx.Preload("Author").Where("id=?", videoId).First(&video)
	author = video.Author
	if actionType == 1 {
		favorite.IsFavorite = true
		// 查看favorite表格是否存在数据，如果已经存在且已经点过赞，则直接返回
		var existingFavorite model.Favorite
		result := tx.Where("user_id = ? AND video_id = ?", userId, videoId).First(&existingFavorite)
		fmt.Println(actionType, "...", existingFavorite.IsFavorite)

		if result.RowsAffected == 0 {
			tx.Create(&favorite)
		} else if existingFavorite.IsFavorite == true {
			return errors.New("请勿重复点赞！")
		} else { //将false更新为true
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Update("is_favorite", true).Error; err != nil {
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
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}
	} else if actionType == 2 {
		//查找是否存在点赞记录，若没有或者已经是没点赞的状态下则无法取消点赞
		var existFavorite model.Favorite
		tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).First(&existFavorite)
		if existFavorite.Id == 0 || existFavorite.IsFavorite == false {
			return errors.New("请勿在没点赞情况下取消点赞")
		}

		fmt.Println(actionType, "...", existFavorite.IsFavorite)
		//将true改为false表示取消点赞
		tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Update("is_favorite", false)

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
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}
	} else {
		return errors.New("无效操作！")
	}

	tx.Commit()
	return nil
}

func GetFavoriteList(sourceId, userId int64) ([]*pb.Video, error) {
	var favorites []model.Favorite
	var resVideoList []*pb.Video
	if err := common.DB.Preload("Video").Preload("Video.Author").Model(&model.Favorite{}).Where("user_id = ?", userId).Find(&favorites).Error; err != nil {
		return nil, err
	}

	//设置视频结构体为为喜欢并且转换为响应结构体
	for _, f := range favorites {
		//检查当前请求用户是否对userId的用户的点赞列表视频点赞
		var isFavorite bool
		isFavorite, err := JudgeFavorite(sourceId, f.VideoId)
		if err != nil {
			return nil, err
		}

		//判断favorite结构体是否点赞
		if f.IsFavorite {
			f.Video.IsFavorite = isFavorite
			resVideoList = append(resVideoList, ToProtoVideo(f.Video))
		}
	}
	return resVideoList, nil
}
