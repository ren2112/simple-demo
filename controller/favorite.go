package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//获取token，actionType，videoId
	token := c.Query("token")
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		response.CommonResp(c, 1, "无效操作！")
		return
	}
	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}

	conn := common.ConnFavoritePool.Get()

	client := pb.NewFavoriteServiceClient(conn)
	resp, err := client.FavoriteAction(c, &pb.DouyinFavoriteActionRequest{ActionType: int32(actionType), Token: token, VideoId: int64(videoID)})
	common.ConnRelationPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	} else {
		response.CommonResp(c, resp.StatusCode, resp.StatusMsg)
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//获取userId，token
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	if err != nil {
		response.CommonResp(c, 1, "用户不存在！")
	}

	conn := common.ConnFavoritePool.Get()

	client := pb.NewFavoriteServiceClient(conn)
	resp, err := client.GetFavoriteList(c, &pb.DouyinFavoriteListRequest{UserId: userId, Token: token})
	common.ConnFavoritePool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	} else {
		response.VideoListResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.VideoList)
	}
}
