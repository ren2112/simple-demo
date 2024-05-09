package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []model.User `json:"user_list"`
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
	case "2": // 取消关注用户
		// 删除关注关系
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

	var userList []model.User
	var followList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("User").Where("follower_user_id=?", userId).Find(&followList)
	for _, v := range followList {
		v.User.IsFollow = true
		userList = append(userList, v.User)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
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
		return
	}

	var userList []model.User
	var followerList []model.Follow
	common.DB.Preload("FollowerUser").Model(&model.Follow{}).Where("user_id=?", userId).Find(&followerList)
	for _, v := range followerList {
		userList = append(userList, v.FollowerUser)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
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

	var userList []model.User
	var followerList []model.Follow
	common.DB.Model(&model.Follow{}).Preload("FollowerUser").Where("user_id=?", userId).Find(&followerList)
	for _, v := range followerList {
		userList = append(userList, v.FollowerUser)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
