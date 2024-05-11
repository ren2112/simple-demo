package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		response.CommonResp(c, 1, "请重新登陆！")
		return
	}
	userId := user.(model.User).Id

	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	actionType := c.Query("action_type")
	var favorite model.Favorite
	favorite.UserId = userId
	favorite.VideoId = int64(videoID)

	//开启事务
	tx := common.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.CommonServerError(c)
		}
	}()

	var author model.User
	var video model.Video

	//获得视频作者的信息，因为要对视频作者的获赞数量更改
	common.DB.Preload("Author").Where("id=?", videoID).First(&video)
	author = video.Author
	if actionType == "1" {
		favorite.IsFavorite = true
		// 查看favorite表格是否存在数据，如果已经存在且已经点过赞，则直接返回
		var existingFavorite model.Favorite
		if result := tx.Where("user_id = ? AND video_id = ?", userId, videoID).First(&existingFavorite); result.RowsAffected == 0 {
			tx.Create(&favorite)
		} else if existingFavorite.IsFavorite == true {
			response.CommonResp(c, 1, "请勿重复点赞")
			return
		} else { //将false更新为true
			tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoID).Update("is_favorite", true)
		}

		// 更新视频的favorite_count字段
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新视频作者的TotalFavorited字段
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新用户的favorite_count字段
		tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))

		response.CommonResp(c, 0, "点赞成功！")
	} else {
		//查找是否存在点赞记录，若没有或者已经是没点赞的状态下则无法取消点赞
		var existFavorite model.Favorite
		tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoID).First(&existFavorite)
		if existFavorite.Id == 0 || existFavorite.IsFavorite == false {
			response.CommonResp(c, 1, "请勿在没点赞情况下取消点赞")
			return
		}

		tx.Model(&model.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoID).Update("is_favorite", false)
		// 更新视频的favorite_count字段
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新视频作者的TotalFavorited字段
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新用户的favorite_count字段
		tx.Model(&model.User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

		response.CommonResp(c, 0, "取消点赞成功")
	}

	tx.Commit()
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")

	var favorites []model.Favorite
	var resVideoList []model.RespVideo
	common.DB.Preload("Video").Preload("Video.Author").Model(&model.Favorite{}).Where("user_id = ?", userId).Find(&favorites)

	//设置视频结构体为为喜欢并且转换为响应结构体
	for _, v := range favorites {
		v.Video.IsFavorite = true
		resVideoList = append(resVideoList, assist.ToRespVideo(v.Video))
	}

	response.VideoListResponseFun(c, response.Response{StatusCode: 0, StatusMsg: "获取成功！"}, resVideoList)
}
