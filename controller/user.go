package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//获取rpc连接
	conn := common.GetUserConnection()

	client := pb.NewUserServiceClient(conn)
	resp, err := client.Regist(c, &pb.DouyinUserRegisterRequest{Username: username, Password: password})
	if err != nil {
		response.CommonResp(c, 1, "注册失败"+err.Error())
		return
	}
	response.UserLoginOk(c, resp.UserId, resp.Token)
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	conn := common.GetUserConnection()

	client := pb.NewUserServiceClient(conn)
	resp, err := client.Login(c, &pb.DouyinUserLoginRequest{Username: username, Password: password})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserLoginOk(c, resp.UserId, resp.Token)
}

func UserInfo(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserLoginRespFail(c, "用户不存在")
	}

	conn := common.GetUserConnection()

	client := pb.NewUserServiceClient(conn)
	resp, err := client.GetUserInfo(c, &pb.DouyinUserRequest{UserId: int64(userId)})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserResponseFun(c, resp.User)
}
