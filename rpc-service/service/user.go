package service

import (
	"context"
	"errors"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (u UserService) Login(ctx context.Context, req *pb.DouyinUserLoginRequest) (*pb.DouyinUserLoginResponse, error) {
	//从redis缓存找
	userPointerFromRedis, err := common.GetCachedUser(ctx, req.Username)
	var user model.User
	if err == nil && userPointerFromRedis != nil {
		user = *userPointerFromRedis
	} else {
		//如果缓存没找到，则用数据库查找是否存在用户
		user = service.GetUserByName(req.Username)
		if user.Id == 0 {
			return &pb.DouyinUserLoginResponse{StatusCode: 1, StatusMsg: "用户名或密码错误！"}, nil
		}
	}

	//校验密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.DouyinUserLoginResponse{StatusCode: 1, StatusMsg: "用户名或密码错误！"}, nil
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.DouyinUserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功！",
		UserId:     user.Id,
		Token:      token,
	}, nil
}

func (u UserService) Regist(ctx context.Context, req *pb.DouyinUserRegisterRequest) (*pb.DouyinUserRegisterResponse, error) {
	if user := service.GetUserByName(req.Username); user.Id != 0 {
		return nil, errors.New("用户已存在！")
	} else {
		//修改为bcrypt加密
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return &pb.DouyinUserRegisterResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
		}
		//创建新用户
		newUser := model.User{
			Name:            req.Username,
			Password:        string(hasedPassword),
			Avatar:          config.DEFAULT_USER_AVATAR_URL,
			BackgroundImage: config.DEFAULT_USER_BG_IMAGE_URL,
			Signature:       config.DEFAULT_USER_BIO,
		}

		//调用服务层数据库操作
		err = service.CreateUser(newUser)
		if err != nil {
			return &pb.DouyinUserRegisterResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
		}

		//获取token
		token, err := common.ReleaseToken(newUser)
		if err != nil {
			return nil, err
		}
		return &pb.DouyinUserRegisterResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功",
			UserId:     user.Id,
			Token:      token,
		}, nil
	}
}

func (u UserService) GetUserInfo(ctx context.Context, req *pb.DouyinUserRequest) (*pb.DouyinUserResponse, error) {
	user := service.GetUserByID(req.UserId)
	if user.Id == 0 {
		return &pb.DouyinUserResponse{StatusCode: 1, StatusMsg: "用户不存在！"}, nil
	} else {
		retUser := service.ToProtoUser(user)
		return &pb.DouyinUserResponse{StatusCode: 0, User: retUser}, nil
	}
}
