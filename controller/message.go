package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil || actionType != "1" {
		response.CommonResp(c, 1, "无效操作！")
	}
	content := c.Query("content")

	conn := common.GetMessageConnection()
	// 建立连接
	client := pb.NewChatServiceClient(conn)
	resp, err := client.ChatAction(c, &pb.DouyinChatActionRequest{Token: token, ToUserId: toUserId, Content: content})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.CommonResp(c, resp.StatusCode, resp.StatusMsg)
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		response.CommonResp(c, 1, "请求失败")
		return
	}
	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, err := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if err != nil {
		response.CommonResp(c, 1, "请求失败")
		return
	}

	conn := common.GetMessageConnection()
	// 建立连接
	client := pb.NewChatServiceClient(conn)
	resp, err := client.GetChatList(c, &pb.DouyinMessageChatRequest{Token: token, ToUserId: toUserId, PreMsgTime: preMsgTime})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.ChatResponseFun(c, response.Response{StatusCode: 0}, resp.MessageList)
}
