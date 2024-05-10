package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	tokenStr := c.Query("token")
	_, claims, _ := common.ParseToken(tokenStr)
	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	actionType := c.Query("action_type")
	var favorite model.Favorite
	favorite.UserId = claims.UserId
	favorite.VideoId = int64(videoID)
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Response{
				StatusCode: 1,
				StatusMsg:  "操作失败！",
			})
		}
	}()

	var author model.User
	var video model.Video
	common.DB.Preload("Author").Where("id=?", videoID).First(&video)
	author = video.Author
	if actionType == "1" {
		// 查看favorite表格是否存在数据，如果已经存在，则不处理
		var existingFavorite model.Favorite
		if result := tx.Where("user_id = ? AND video_id = ?", claims.UserId, videoID).First(&existingFavorite); result.RowsAffected == 0 {
			tx.Create(&favorite)
		}

		// 更新UserVideo表，表示用户对该视频点赞
		var userVideo model.UserVideo
		if result := tx.Where("user_id = ? AND video_id = ?", claims.UserId, videoID).First(&userVideo); result.RowsAffected == 0 {
			userVideo = model.UserVideo{
				UserID:     claims.UserId,
				VideoID:    int64(videoID),
				IsFavorite: true,
			}
			tx.Create(&userVideo)
		} else {
			tx.Model(&userVideo).Where("user_id = ? AND video_id = ?", userVideo.UserID, userVideo.VideoID).Update("is_favorite", true)
		}

		// 更新视频的favorite_count字段
		tx.Model(&model.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))

		// 更新视频作者的TotalFavorited字段
		tx.Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))

		// 更新用户的favorite_count字段
		tx.Model(&model.User{}).Where("id = ?", claims.UserId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "点赞成功",
		})
	} else {
		tx.Where("user_id = ? AND video_id = ?", claims.UserId, videoID).Delete(&model.Favorite{})
		// 更新UserVideo表，表示用户取消对该视频的点赞
		tx.Model(&model.UserVideo{}).Where("user_id = ? AND video_id = ?", claims.UserId, videoID).Update("is_favorite", false)

		// 更新视频的favorite_count字段
		tx.Model(&model.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

		// 更新视频作者的TotalFavorited字段
		tx.Model(&model.User{}).Where("id = ?", author.Id).Update("total_favorited", gorm.Expr("total_favorited - ?", 1))

		// 更新用户的favorite_count字段
		tx.Model(&model.User{}).Where("id = ?", claims.UserId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取消点赞成功",
		})
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取成功！",
		},
		VideoList: resVideoList,
	})
}
