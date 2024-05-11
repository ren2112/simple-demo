package response

import (
	"github.com/RaymondCode/simple-demo/model"
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
	User model.RespUser `json:"user"`
}

type UserListResponse struct {
	Response
	UserList []model.RespUser `json:"user_list"`
}

func UserLoginRespFail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 1, StatusMsg: msg},
	})
}

func UserLoginOk(c *gin.Context, Id int64, token string) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   Id,
		Token:    token,
	})
}

func UserResponseFun(c *gin.Context, respUser model.RespUser) {
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     respUser,
	})
}

func UserListResponseFun(c *gin.Context, response Response, respUserList []model.RespUser) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: response,
		UserList: respUserList,
	})
}
