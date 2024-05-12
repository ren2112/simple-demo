package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func CommentAction(actionType string, user model.User, videoId int64, text string, commentId string, respComment *model.RespComment) error {
	userId := user.Id

	//开始对comment表格操作
	var comment model.Comment
	comment.VideoId = videoId
	comment.User = user
	comment.UserId = userId

	// 添加评论
	if actionType == "1" {
		// 开启事务
		tx := common.DB.Begin()
		comment.Content = text
		tx.Create(&comment)

		// 增加评论数
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}

		tx.Commit()

		respComment = &model.RespComment{
			Id:         comment.Id,
			User:       ToRespUser(user),
			Content:    text,
			CreateDate: comment.CreatedAt.Format("2006-01-02 15:04"),
		}
		return nil
	} else {
		// 开启事务
		tx := common.DB.Begin()
		// 删除评论
		if err := tx.Where("id=?", commentId).Delete(&comment).Error; err != nil {
			return err
		}

		// 减少评论数
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	}
	return nil
}

func GetCommentList(videoId string) ([]model.RespComment, error) {
	var comments []model.Comment
	var respComments []model.RespComment
	if err := common.DB.Preload("User").Model(&model.Comment{}).Where("video_id=?", videoId).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	//转化为响应结构体
	for _, v := range comments {
		respComments = append(respComments, ToRespComment(v))
	}
	return respComments, nil
}

func ToRespComment(comment model.Comment) model.RespComment { // 转换 Comment 中的 User 字段为 RespUser
	respComment := model.RespComment{
		Id:         comment.Id,
		User:       ToRespUser(comment.User),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化时间,
	}

	return respComment
}
