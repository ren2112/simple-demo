package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if user := service.GetUserByName(username); user.Id != 0 {
		response.UserLoginRespFail(c, "用户已存在")
		return
	} else {
		//修改为bcrypt加密
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.UserLoginRespFail(c, "注册失败")
			return
		}
		//创建新用户
		newUser := model.User{
			Name:            username,
			Password:        string(hasedPassword),
			Avatar:          config.DEFAULT_USER_AVATAR_URL,
			BackgroundImage: config.DEFAULT_USER_BG_IMAGE_URL,
			Signature:       config.DEFAULT_USER_BIO,
		}

		//调用服务层数据库操作
		err = service.CreateUser(newUser)
		if err != nil {
			response.UserLoginRespFail(c, "注册失败！")
			return
		}

		//获取token
		token, err := common.ReleaseToken(newUser)
		if err != nil {
			response.UserLoginRespFail(c, "发送token失败")
			return
		}
		response.UserLoginOk(c, newUser.Id, token)
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//查找是否存在用户
	user := service.GetUserByName(username)
	if user.Id == 0 {
		response.UserLoginRespFail(c, "用户名或者密码错误")
		return
	}

	//校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.UserLoginRespFail(c, "用户名或者密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.UserLoginRespFail(c, "发送token失败")
		return
	}
	response.UserLoginOk(c, user.Id, token)
}

func UserInfo(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserLoginRespFail(c, "用户不存在")
	}
	user := service.GetUserByID(int64(userId))
	if user.Id == 0 {
		response.UserLoginRespFail(c, "用户不存在")
	} else {
		response.UserResponseFun(c, service.ToRespUser(user))
	}
}
