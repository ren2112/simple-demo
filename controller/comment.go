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
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			response.CommonResp(c, 1, err.Error())
			return
		}

		tx.Commit()

		respComment := model.RespComment{
			Id:         comment.Id,
			User:       assist.ToRespUser(user),
			Content:    text,
			CreateDate: comment.CreatedAt.Format("2006-01-02 15:04"),
		}
		response.CommentActionResponseFun(c, response.Response{StatusCode: 0}, respComment)
	} else {
		// 开启事务
		tx := common.DB.Begin()
		// 删除评论
		tx.Where("id=?", commentId).Delete(&comment)

		// 减少评论数
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			response.CommonResp(c, 1, err.Error())
			return
		}

		tx.Commit()
		response.CommonResp(c, 0, "")
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
	response.CommentListResponseFun(c, response.Response{StatusCode: 0}, respComments)
}
