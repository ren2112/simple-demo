package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	//token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	if content == "" {
		response.CommonResp(c, 1, "消息不能为空")
		return
	}

	if user, exist := c.Get("user"); exist {
		//添加信息于数据库
		userIdTarget, _ := strconv.Atoi(toUserId)
		var message model.Message
		message.Content = content
		message.FromUserId = user.(model.User).Id
		message.ToUserId = int64(userIdTarget)
		common.DB.Create(&message)
		response.CommonResp(c, 0, "")
	} else {
		response.CommonResp(c, 1, "用户不存在！")
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	//token := c.Query("token")
	toUserId := c.Query("to_user_id")
	preMsgTimeStr := c.Query("pre_msg_time")

	if user, exist := c.Get("user"); exist {
		var preMsgTime int64
		var err error
		preMsgTime, err = strconv.ParseInt(preMsgTimeStr, 10, 64)
		if err != nil {
			response.CommonResp(c, 1, "请求失败")
			return
		}
		preMsgTimeUTC := time.Unix(0, preMsgTime*int64(time.Millisecond))

		userIdTarget, _ := strconv.Atoi(toUserId)
		fromUserId := user.(model.User).Id
		var messageList []model.Message
		common.DB.Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)", userIdTarget, fromUserId, fromUserId, userIdTarget).
			Where("created_at > ?", preMsgTimeUTC).
			Order("created_at").
			Find(&messageList)

		var resMessageList []model.RespMessage
		for _, v := range messageList {
			var resMessage = assist.ToRespMessage(v)
			resMessageList = append(resMessageList, resMessage)
		}
		response.ChatResponseFun(c, response.Response{StatusCode: 0}, resMessageList)
	} else {
		response.CommonResp(c, 1, "用户不存在！")
	}
}
