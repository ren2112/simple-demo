package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
)

func CreateUser(newUser model.User) error {
	err := common.DB.Create(&newUser).Error
	return err
}

func GetUserByName(userName string) (user model.User) {
	common.DB.Where("name=?", userName).First(&user)
	return user
}

func GetUserByID(ID int64) (user model.User) {
	common.DB.Find(&user, ID)
	return user
}

func ToRespUser(user model.User) model.RespUser {
	return model.RespUser{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		WorkCount:       user.WorkCount,
		TotalFavorited:  user.TotalFavorited,
		FavoriteCount:   user.FavoriteCount,
	}
}

func ToProtoUser(user model.User) pb.User {
	return pb.User{
		Id:              user.Id,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		WorkCount:       user.WorkCount,
		TotalFavorited:  user.TotalFavorited,
		FavoriteCount:   user.FavoriteCount,
	}
}
