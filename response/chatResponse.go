package response

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatResponse struct {
	Response
	MessageList []model.RespMessage `json:"message_list"`
}

func ChatResponseFun(c *gin.Context, response Response, resMessageList []model.RespMessage) {
	c.JSON(http.StatusOK, ChatResponse{Response: response, MessageList: resMessageList})
}
