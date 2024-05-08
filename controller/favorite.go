package controller

import (
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
	videoId, err := strconv.Atoi(c.Query("video_id"))
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
	favorite.VideoId = int64(videoId)

	if actionType == "1" {
		//查看favorite表格是否存在数据，如果已经存在，则不处理
		var existingFavorite model.Favorite
		result := common.DB.Where("user_id = ? AND video_id = ?", claims.UserId, videoId).First(&existingFavorite)
		if result.RowsAffected == 0 {
			common.DB.Create(&favorite)
		}
		// 开启事务,对点赞数量更新的时候需要上行级锁
		tx := common.DB.Begin()
		var userVideo model.UserVideo
		result = tx.Where("user_id = ? AND video_id = ?", claims.UserId, videoId).First(&userVideo)
		if result.RowsAffected == 0 {
			// 如果关联记录不存在，则创建新记录
			userVideo = model.UserVideo{
				UserID:     claims.UserId,
				VideoID:    int64(videoId),
				IsFavorite: true,
			}
			tx.Create(&userVideo)
		} else {
			// 如果关联记录存在，则更新isfavorite为true
			tx.Model(&userVideo).Where("user_id = ? AND video_id = ?", userVideo.UserID, userVideo.VideoID).Update("is_favorite", true)
		}
		// 获取行级锁
		tx.Exec("SELECT id FROM videos WHERE id = ? FOR UPDATE", videoId)
		// 执行更新操作
		tx.Model(&model.Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		// 提交事务
		tx.Commit()

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "点赞成功",
		})
	} else {
		common.DB.Where("user_id = ? AND video_id = ?", claims.UserId, videoId).Delete(&model.Favorite{})
		// 开启事务，对点赞数量更新的时候需要上行级锁
		tx := common.DB.Begin()

		userVideo := model.UserVideo{
			UserID:  claims.UserId,
			VideoID: int64(videoId),
		}
		tx.Model(&userVideo).Where("user_id = ? AND video_id = ?", userVideo.UserID, userVideo.VideoID).Update("is_favorite", false)
		// 获取行级锁
		tx.Exec("SELECT id FROM videos WHERE id = ? FOR UPDATE", videoId)
		// 执行更新操作
		tx.Model(&model.Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		// 提交事务
		tx.Commit()

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取消点赞成功",
		})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")
	var favorites []model.Favorite
	var videoList []model.Video
	common.DB.Preload("Video").Preload("Video.Author").Model(&model.Favorite{}).Where("user_id = ?", userId).Find(&favorites)
	for _, v := range favorites {
		v.Video.IsFavorite = true
		videoList = append(videoList, v.Video)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取成功！",
		},
		VideoList: videoList,
	})
}
