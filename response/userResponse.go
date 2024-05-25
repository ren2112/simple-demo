package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User *pb.User `json:"user"`
}

type UserListResponse struct {
	Response
	UserList []*pb.User `json:"user_list"`
}

func UserLoginRespFail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 1, StatusMsg: msg},
	})
}

func UserLoginResp(c *gin.Context, Id int64, token string, statusCode int32, statusMsg string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: statusCode, StatusMsg: statusMsg},
		UserId:   Id,
		Token:    token,
	})
}

func UserResponseFun(c *gin.Context, respUser *pb.User) {
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     respUser,
	})
}

func UserListResponseFun(c *gin.Context, response Response, respUserList []*pb.User) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: response,
		UserList: respUserList,
	})
}
