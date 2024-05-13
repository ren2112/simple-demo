package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

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
		response.CommonResp(c, 1, "请求视频失败")
		return
	}
	var videoList = []model.Video{}
	latestTimeUTC := time.Unix(0, latestTime*int64(time.Millisecond))

	videoList, err = service.FeedVideoList(latestTimeUTC)
	if err != nil {
		response.CommonResp(c, 1, "请求视频失败")
		return
	}

	//若用户登录了，需要判断视频作者是否关注以及是否对视频点赞
	if claims.UserId != 0 {
		for i, v := range videoList {
			//查找是否点赞
			var isFavorite bool
			isFavorite, err = service.JudgeFavorite(claims.UserId, v.Id)
			if err != nil {
				response.CommonResp(c, 1, "点赞数据请求失败！")
			}
			videoList[i].IsFavorite = isFavorite

			//	查找videoList的author里面is_follow
			// 查找作者是否被当前用户关注
			//var follow model.Follow
			//err = common.DB.Where("user_id = ? AND follower_user_id = ?", v.Author.Id, claims.UserId).First(&follow).Error
			err = service.JudgeRelation(v.AuthorID, claims.UserId)
			if err == nil {
				videoList[i].Author.IsFollow = true
			} else if err == gorm.ErrRecordNotFound {
				videoList[i].Author.IsFollow = false
			} else {
				response.CommonResp(c, 1, "获取视频失败，请重试！")
				return
			}
		}
	}

	//将videoList的每个元素赋值给respvideoList
	var respVideoList []model.RespVideo
	for _, v := range videoList {
		respVideo := service.ToRespVideo(v)
		//补全视频url
		respVideo.PlayUrl = config.SERVER_RESOURCES + v.PlayUrl
		respVideo.CoverUrl = config.SERVER_RESOURCES + v.CoverUrl
		respVideoList = append(respVideoList, respVideo)
	}

	//获取视频流里最早视频时间
	var responseTime int64
	if len(videoList) == 0 {
		responseTime = time.Now().Unix()
	} else {
		responseTime = videoList[len(videoList)-1].CreatedAt.Unix()
	}
	response.FeedResponseFun(c, respVideoList, responseTime)
}
