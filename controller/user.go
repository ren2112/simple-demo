package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if user := assist.GetUserByName(username); user.Id != 0 {
		response.UserLoginRespFail(c, "用户已存在")
		return
	} else {
		//修改为bcrypt加密
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.UserLoginRespFail(c, "注册失败")
			return
		}
		newUser := model.User{
			Name:     username,
			Password: string(hasedPassword),
		}
		common.DB.Create(&newUser)
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
	user := assist.GetUserByName(username)
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
	user, ok := c.Get("user")
	if !ok {
		response.UserLoginRespFail(c, "用户不存在")
	} else {
		response.UserResponseFun(c, assist.ToRespUser(user.(model.User)))
	}
}
