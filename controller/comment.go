package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	commentId := c.Query("comment_id")
	text := c.Query("comment_text")

	//获取userid
	user, ok := c.Get("user")
	if !ok {
		response.CommonResp(c, 1, "请先登录！")
	}

	respComment := model.RespComment{}
	err := service.CommentAction(actionType, user.(model.User), int64(videoId), text, commentId, &respComment)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	}
	response.CommentActionResponseFun(c, response.Response{StatusCode: 0}, respComment)
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")

	respComments, err := service.GetCommentList(videoId)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	}
	response.CommentListResponseFun(c, response.Response{StatusCode: 0}, respComments)
}
