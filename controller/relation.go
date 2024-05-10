package controller

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []model.RespUser `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
		return
	}

	targetID, err := strconv.Atoi(c.Query("to_user_id"))
	actionType := c.Query("action_type")
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid target user ID"})
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
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不能关注自己！"})
			return
		}

		// 检查是否已存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id = ? AND follower_user_id = ?", follow.UserId, follow.FollowerUserId).First(&existingFollow)
		if result.RowsAffected == 0 {
			// 不存在则创建关注关系
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败！"})
				return
			}

			// 更新关注者的关注数加一
			if err := tx.Model(&model.User{}).Where("id = ?", user.(model.User).Id).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
				return
			}

			// 更新被关注者的粉丝数加一
			if err := tx.Model(&model.User{}).Where("id = ?", targetID).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
				return
			}
			//	如果已经存在告知不可重复关注
		} else {
			if err := tx.Create(&follow).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请勿重复关注！"})
				return
			}
		}

	case "2": // 取消关注用户
		// 删除关注关系，首先判断是否存在关注关系
		var existingFollow model.Follow
		result := tx.Where("user_id=? and follower_user_id=?", targetID, user.(model.User).Id).First(&existingFollow)
		if result.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "请勿重复取消关注！"})
			return
		}

		if err := tx.Where("user_id = ? AND follower_user_id = ?", targetID, user.(model.User).Id).Delete(&model.Follow{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			return
		}

		// 更新关注者的关注数减一
		if err := tx.Model(&model.User{}).Where("id = ?", user.(model.User).Id).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			return
		}

		// 更新被关注者的粉丝数减一
		if err := tx.Model(&model.User{}).Where("id = ?", targetID).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
			return
		}
	default:
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作失败"})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "操作成功！"})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
			},
			UserList: nil,
		})
	}

	var respUserList []model.RespUser
	var followList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("User").Where("follower_user_id=?", userId).Find(&followList)
	//转换为响应结构体
	for _, v := range followList {
		v.User.IsFollow = true
		respUserList = append(respUserList, assist.ToRespUser(v.User))
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: respUserList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
			},
			UserList: nil,
		})
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	sourceUser, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请先登录！",
			},
			UserList: nil,
		})
	}
	sourceId := sourceUser.(model.User).Id

	var respUserList []model.RespUser
	var followerList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList)
	//转换结构体
	for _, v := range followerList {
		//注意，这里需要传递请求发起者的id，而不是user_id，因为user_id是被查看粉丝列表的人
		v.FollowerUser.IsFollow = assist.IsFollowed(sourceId, v.FollowerUserId)
		respUserList = append(respUserList, assist.ToRespUser(v.FollowerUser))
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: respUserList,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
			},
			UserList: nil,
		})
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	sourceUser, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请先登录！",
			},
			UserList: nil,
		})
	}
	sourceId := sourceUser.(model.User).Id

	var respUserList []model.RespUser
	var followerList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList)
	//转换结构体
	for _, v := range followerList {
		//注意，这里需要传递请求发起者的id，而不是user_id，因为user_id是被查看粉丝列表的人
		v.FollowerUser.IsFollow = assist.IsFollowed(sourceId, v.FollowerUserId)
		respUserList = append(respUserList, assist.ToRespUser(v.FollowerUser))
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: respUserList,
	})
}
