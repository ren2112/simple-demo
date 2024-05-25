package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"gorm.io/gorm"
)

func RelationAction(actionType string, userId int64, targetId int) error {
	// 开启事务
	tx := common.DB.Begin()

	switch actionType {
	case "1": // 关注用户
		// 创建关注关系
		var follow model.Follow
		follow.FollowerUserId = userId
		follow.UserId = int64(targetId)
		if follow.UserId == follow.FollowerUserId {
			tx.Rollback()
			return errors.New("不能关注自己！")
		}

		// 检查是否已存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id = ? AND follower_user_id = ?", follow.UserId, follow.FollowerUserId).First(&existingFollow)
		if result.RowsAffected == 0 {
			// 不存在则创建关注关系
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				return errors.New("操作失败！")
			}

			// 更新关注者的关注数加一
			if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				return errors.New("操作失败！")
			}

			// 更新被关注者的粉丝数加一
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", targetId).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				return errors.New("操作失败！")
			}
			//	如果已经存在告知不可重复关注
		} else {
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				return errors.New("请勿重复关注！")
			}
		}

	case "2": // 取消关注用户
		// 删除关注关系，首先判断是否存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id=? and follower_user_id=?", targetId, userId).First(&existingFollow)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return errors.New("请勿重复取消关注！")
		}

		if err := tx.Where("user_id = ? AND follower_user_id = ?", targetId, userId).Delete(&model.Follow{}).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新关注者的关注数减一
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}

		// 更新被关注者的粉丝数减一
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.User{}).Where("id = ?", targetId).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			return errors.New("操作失败！")
		}
	default:
		tx.Rollback()
		return errors.New("操作失败！")
	}

	// 提交事务
	tx.Commit()
	return nil
}

func GetFollowList(sourceId int64, userId int) ([]*pb.User, error) {
	var respUserList []*pb.User
	var followList []model.Follow
	if err := common.DB.Model(&model.Follow{}).Preload("User").Where("follower_user_id=?", userId).Find(&followList).Error; err != nil {
		return nil, err
	}
	//判断发起请求的用户是否关注了关注列表的用户，并且转换为响应结构体
	for _, v := range followList {
		if sourceId != 0 {
			v.User.IsFollow = IsFollowed(sourceId, v.UserId)
		}
		respUserList = append(respUserList, ToProtoUser(v.User))
	}
	return respUserList, nil
}

func GetFollowerList(sourceId int64, userId int64) ([]*pb.User, error) {
	var respUserList []*pb.User
	var followerList []model.Follow
	if err := common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList).Error; err != nil {
		return nil, err
	}
	//转换结构体
	for _, v := range followerList {
		//注意，这里需要传递请求发起者的id，而不是user_id，因为user_id是被查看粉丝列表的人
		if sourceId != 0 {
			v.FollowerUser.IsFollow = IsFollowed(sourceId, v.FollowerUserId)
		}
		respUserList = append(respUserList, ToProtoUser(v.FollowerUser))
	}
	return respUserList, nil
}

func IsFollowed(aUserId, bUserId int64) bool {
	// 查询是否a关注了指定b用户的逻辑
	follow := model.Follow{}
	common.DB.Where("user_id=? and follower_user_id=?", bUserId, aUserId).Find(&follow)
	if follow.Id != 0 {
		return true
	}
	return false
}

func ToProtoFriend(user *pb.User, msg string, msgType int64) *pb.FriendUser {
	return &pb.FriendUser{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
		Message:         msg,
		MsgType:         msgType,
	}
}
