package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	_, claims, _ := common.ParseToken(c.Query("token"))

	latestTimeStr := c.Query("latest_time")
	var latestTime int64
	var err error
	if latestTimeStr == "" {
		latestTime = time.Now().Unix() * 1000
	} else {
		latestTime, err = strconv.ParseInt(latestTimeStr, 10, 64)
	}
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频失败！请重试" + c.Query("latest_time")},
		})
		return
	}
	var videoList = []model.Video{}
	latestTimeUTC := time.Unix(0, latestTime*int64(time.Millisecond))
	common.DB.Preload("Author").Model(&model.Video{}).Where("created_at < ?", latestTimeUTC).Find(&videoList)

	for i, v := range videoList {
		var userVideo model.UserVideo
		//查找是否点赞
		common.DB.Where("user_id=? and video_id=?", claims.UserId, v.Id).First(&userVideo)
		videoList[i].IsFavorite = userVideo.IsFavorite

		//	查找videoList的author里面is_follow
		// 查找作者是否被当前用户关注
		var follow model.Follow
		err := common.DB.Where("user_id = ? AND follower_user_id = ?", v.Author.Id, claims.UserId).First(&follow).Error
		if err == nil {
			videoList[i].Author.IsFollow = true
		} else if err == gorm.ErrRecordNotFound {
			videoList[i].Author.IsFollow = false
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "获取视频失败！请重试" + c.Query("latest_time")},
			})
			return
		}
	}

	var responseTime int64
	if len(videoList) == 0 {
		responseTime = time.Now().Unix()
	} else {
		responseTime = videoList[len(videoList)-1].CreatedAt.Unix()
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  responseTime * 1000,
	})
}
