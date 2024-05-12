package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
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
	err = service.FavoriteAction(actionType, userId, int64(videoID))
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	}
	response.CommonResp(c, 0, "操作成功！")
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")

	resVideoList, err := service.GetFavoriteList(userId)
	if err != nil {
		response.CommonResp(c, 1, "获取点赞视频异常！")
	}
	response.VideoListResponseFun(c, response.Response{StatusCode: 0, StatusMsg: "获取成功！"}, resVideoList)
}
