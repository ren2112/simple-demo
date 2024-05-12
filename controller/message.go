package controller

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
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
		err := service.CreateMessage(toUserId, content, user.(model.User).Id)
		if err != nil {
			response.CommonResp(c, 1, "发送消息失败请重试")
		}
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

		//处理时间
		preMsgTime, err = strconv.ParseInt(preMsgTimeStr, 10, 64)
		if err != nil {
			response.CommonResp(c, 1, "请求失败")
			return
		}
		preMsgTimeUTC := time.Unix(0, preMsgTime*int64(time.Millisecond))

		fromUserId := user.(model.User).Id
		var resMessageList []model.RespMessage
		resMessageList, err = service.GetMessageList(toUserId, fromUserId, preMsgTimeUTC)
		if err != nil {
			response.CommonResp(c, 1, "获取聊天列表失败！")
		} else {
			response.ChatResponseFun(c, response.Response{StatusCode: 0}, resMessageList)
		}
	} else {
		response.CommonResp(c, 1, "用户不存在！")
	}
}
