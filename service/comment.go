package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	redislock "github.com/jefferyjob/go-redislock"
	"gorm.io/gorm"
)

func CommentAction(actionType int32, user *pb.User, videoId int64, text string, commentId int64, respComment *pb.Comment) error {
	ctx := context.Background()
	userId := user.Id

	//开始对comment表格操作
	var comment model.Comment
	comment.VideoId = videoId
	comment.UserId = userId

	//增加评论数需要上分布式锁，保证数据一致
	lock := redislock.New(ctx, common.RedisClient, fmt.Sprintf("comment_video:%d", videoId), redislock.WithAutoRenew())
	err := lock.Lock()
	if err != nil {
		return errors.New("操作失败！")
	}
	defer lock.UnLock()

	// 添加评论
	if actionType == 1 {
		// 开启事务
		tx := common.DB.Begin()
		comment.Content = text
		tx.Create(&comment)

		// 增加评论数
		result := tx.Set("gorm:query_option", "FOR UPDATE").
			Model(&model.Video{}).
			Where("id = ?", videoId).
			Update("comment_count", gorm.Expr("comment_count + ?", 1))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("视频不存在！")
		}

		tx.Commit()

		respComment.Id = comment.Id
		respComment.User = user
		respComment.Content = text
		respComment.CreateDate = comment.CreatedAt.Format(config.DATETIME_FORMAT)

		return nil
	} else if actionType == 2 {
		// 开启事务
		tx := common.DB.Begin()
		// 删除评论
		if err := tx.Where("id=?", commentId).Delete(&comment).Error; err != nil {
			return err
		}

		// 减少评论数
		if result := tx.Set("gorm:query_option", "FOR UPDATE").
			Model(&model.Video{}).
			Where("id = ?", videoId).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)); result.Error != nil {

			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return errors.New("视频不存在")
		}

		tx.Commit()
	} else {
		return errors.New("无效操作！")
	}
	return nil
}

func GetCommentList(videoId int64) ([]*pb.Comment, error) {
	var comments []model.Comment
	var respComments []*pb.Comment
	if err := common.DB.Preload("User").Model(&model.Comment{}).Where("video_id=?", videoId).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	//转化为响应结构体
	for _, v := range comments {
		respComments = append(respComments, ToProtoComment(v))
	}
	return respComments, nil
}

func ToProtoComment(comment model.Comment) *pb.Comment { // 转换 Comment 中的 User 字段为 RespUser
	respComment := pb.Comment{
		Id:         comment.Id,
		User:       ToProtoUser(comment.User),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化时间,
	}
	return &respComment
}
