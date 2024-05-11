package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		response.CommonResp(c, 1, "请重新登陆")
		return
	}

	targetID, err := strconv.Atoi(c.Query("to_user_id"))
	actionType := c.Query("action_type")
	if err != nil {
		response.CommonResp(c, 1, "用户不存在")
		return
	}

	// 开启事务
	tx := common.DB.Begin()

	switch actionType {
	case "1": // 关注用户
		// 创建关注关系
		var follow model.Follow
		follow.FollowerUserId = user.(model.User).Id
		follow.UserId = int64(targetID)
		if follow.UserId == follow.FollowerUserId {
			tx.Rollback()
			response.CommonResp(c, 1, "不能关注自己")
			return
		}

		// 检查是否已存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id = ? AND follower_user_id = ?", follow.UserId, follow.FollowerUserId).First(&existingFollow)
		if result.RowsAffected == 0 {
			// 不存在则创建关注关系
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				response.CommonResp(c, 1, "操作失败")
				return
			}

			// 更新关注者的关注数加一
			if err := tx.Model(&model.User{}).Where("id = ?", user.(model.User).Id).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				response.CommonResp(c, 1, "操作失败")
				return
			}

			// 更新被关注者的粉丝数加一
			if err := tx.Model(&model.User{}).Where("id = ?", targetID).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				response.CommonResp(c, 1, "操作失败")
				return
			}
			//	如果已经存在告知不可重复关注
		} else {
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				response.CommonResp(c, 1, "请勿重复关注")
				return
			}
		}

	case "2": // 取消关注用户
		// 删除关注关系，首先判断是否存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id=? and follower_user_id=?", targetID, user.(model.User).Id).First(&existingFollow)
		if result.RowsAffected == 0 {
			tx.Rollback()
			response.CommonResp(c, 1, "请勿重复取消关注")
			return
		}

		if err := tx.Where("user_id = ? AND follower_user_id = ?", targetID, user.(model.User).Id).Delete(&model.Follow{}).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新关注者的关注数减一
		if err := tx.Model(&model.User{}).Where("id = ?", user.(model.User).Id).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}

		// 更新被关注者的粉丝数减一
		if err := tx.Model(&model.User{}).Where("id = ?", targetID).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, "操作失败")
			return
		}
	default:
		tx.Rollback()
		response.CommonResp(c, 1, "操作失败")
		return
	}

	// 提交事务
	tx.Commit()

	response.CommonResp(c, 0, "操作成功")
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败"}, nil)
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的关注列表里面人的是否关注
	sourceUser, ok := c.Get("user")
	if !ok {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	}
	sourceId := sourceUser.(model.User).Id

	var respUserList []model.RespUser
	var followList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("User").Where("follower_user_id=?", userId).Find(&followList)
	//判断发起请求的用户是否关注了关注列表的用户，并且转换为响应结构体
	for _, v := range followList {
		if sourceId != 0 {
			v.User.IsFollow = assist.IsFollowed(sourceId, v.UserId)
		}
		respUserList = append(respUserList, assist.ToRespUser(v.User))
	}

	response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	sourceUser, ok := c.Get("user")
	if !ok {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	}
	sourceId := sourceUser.(model.User).Id

	var respUserList []model.RespUser
	var followerList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList)
	//转换结构体
	for _, v := range followerList {
		//注意，这里需要传递请求发起者的id，而不是user_id，因为user_id是被查看粉丝列表的人
		if sourceId != 0 {
			v.FollowerUser.IsFollow = assist.IsFollowed(sourceId, v.FollowerUserId)
		}
		respUserList = append(respUserList, assist.ToRespUser(v.FollowerUser))
	}

	response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	sourceUser, ok := c.Get("user")
	if !ok {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	}
	sourceId := sourceUser.(model.User).Id

	var respUserList []model.RespUser
	var followerList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList)
	//转换结构体
	for _, v := range followerList {
		//注意，这里需要传递请求发起者的id，而不是user_id，因为user_id是被查看粉丝列表的人
		if sourceId != 0 {
			v.FollowerUser.IsFollow = assist.IsFollowed(sourceId, v.FollowerUserId)
		}
		respUserList = append(respUserList, assist.ToRespUser(v.FollowerUser))
	}

	response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
}
