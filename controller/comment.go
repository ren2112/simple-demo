package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	//获取6个参数：videoId,actionType,commentId,text,token,user
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		response.CommonResp(c, 1, "无效操作！")
		return
	}
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		response.CommonResp(c, 1, "无效操作！")
		return
	}
	commentIdStr := c.Query("comment_id")
	var commentId int64
	if commentIdStr != "" {
		commentId, err = strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			response.CommonResp(c, 1, "无效操作！")
			return
		}
	}
	text := c.Query("comment_text")
	token := c.Query("token")
	user, ok := c.Get("user")
	if !ok {
		response.CommonResp(c, 1, "请先登录！")
		return
	}

	conn := common.ConnCommentPool.Get()

	client := pb.NewCommentServiceClient(conn)
	resp, err := client.CommentAction(c, &pb.DouyinCommentActionRequest{Token: token, VideoId: videoId, ActionType: int32(actionType), CommentText: text, CommentId: commentId, User: service.ToProtoUser(user.(model.User))})
	common.ConnCommentPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.CommentActionResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.Comment)
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		response.CommonResp(c, 1, "无效操作！")
	}

	conn := common.ConnCommentPool.Get()

	client := pb.NewCommentServiceClient(conn)
	resp, err := client.GetCommentList(c, &pb.DouyinCommentListRequest{VideoId: videoId})
	common.ConnCommentPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.CommentListResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.CommentList)
}
