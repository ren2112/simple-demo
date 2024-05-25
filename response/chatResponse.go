package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatResponse struct {
	Response
	MessageList []*pb.Message `json:"message_list"`
}

func ChatResponseFun(c *gin.Context, response Response, resMessageList []*pb.Message) {
	c.JSON(http.StatusOK, ChatResponse{Response: response, MessageList: resMessageList})
}
