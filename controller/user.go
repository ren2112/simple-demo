package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if user := assist.GetUserByName(username); user.Id != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
		return
	} else {
		salt := fmt.Sprintf("%06d", rand.Int31())
		password = utils.MakePassword(password, salt)
		newUser := model.User{
			Name:     username,
			Password: password,
			Salt:     salt,
		}
		common.DB.Create(&newUser)
		token, err := common.ReleaseToken(newUser)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "发放token失败"},
			})
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := assist.GetUserByName(username)
	if user.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码错误"},
		})
		return
	}
	if utils.ValidPassword(password, user.Salt, user.Password) == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或者密码错误"},
		})
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "发放token失败"},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	tokenStr := c.Query("token")

	_, claims, _ := common.ParseToken(tokenStr)
	if user := assist.GetUserByID(claims.UserId); user.Id != 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
}
