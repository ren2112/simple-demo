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

type CommentListResponse struct {
	Response
	CommentList []model.RespComment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment model.RespComment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	tokenStr := c.Query("token")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	commentId := c.Query("comment_id")

	//获取userid
	_, claims, _ := common.ParseToken(tokenStr)
	userId := claims.UserId
	user := assist.GetUserByID(userId)

	//开始对comment表格操作
	var comment model.Comment
	comment.VideoId = int64(videoId)
	comment.User = user
	comment.UserId = userId

	// 添加评论
	if actionType == "1" {
		// 开启事务
		tx := common.DB.Begin()
		text := c.Query("comment_text")
		comment.Content = text
		tx.Create(&comment)

		// 增加评论数
		tx.Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1))

		tx.Commit()

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment: model.RespComment{
				Id:         comment.Id,
				User:       assist.ToRespUser(user),
				Content:    text,
				CreateDate: comment.CreatedAt.Format("2006-01-02 15:04"),
			},
		})
		return
	} else {
		// 开启事务
		tx := common.DB.Begin()
		// 删除评论
		tx.Where("id=?", commentId).Delete(&comment)

		// 减少评论数
		tx.Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))

		tx.Commit()
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")

	var comments []model.Comment
	var respComments []model.RespComment
	common.DB.Preload("User").Model(&model.Comment{}).Where("video_id=?", videoId).Order("created_at DESC").Find(&comments)

	//转化为响应结构体
	for _, v := range comments {
		respComments = append(respComments, assist.ToRespComment(v))
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: respComments,
	})
}
