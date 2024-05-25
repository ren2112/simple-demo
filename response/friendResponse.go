package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FriendResponse struct {
	Response
	UserList []*pb.FriendUser `json:"user_list"`
}

func FriendListResponseFun(c *gin.Context, response Response, respUserList []*pb.FriendUser) {
	c.JSON(http.StatusOK, FriendResponse{
		Response: response,
		UserList: respUserList,
	})
}
