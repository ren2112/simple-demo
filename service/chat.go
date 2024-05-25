package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CreateMessage(toUserId int64, content string, userId int64) error {
	var message model.Message
	message.Content = content
	message.FromUserId = userId
	message.ToUserId = toUserId
	if err := common.DB.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

func GetMessageList(userIdTarget int64, fromUserId int64, preMsgTimeUTC time.Time) ([]*pb.Message, error) {
	var messageList []model.Message
	if err := common.DB.Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)", userIdTarget, fromUserId, fromUserId, userIdTarget).
		Where("created_at > ?", preMsgTimeUTC).
		Order("created_at").
		Find(&messageList).Error; err != nil {
		return nil, err
	}

	var resMessageList []*pb.Message
	for _, v := range messageList {
		var resMessage = ToProtoMessage(v)
		resMessageList = append(resMessageList, resMessage)
	}
	return resMessageList, nil
}

func GetLatestMessage(toUserId, sourceId int64) (*model.Message, error) {
	var message model.Message
	// 根据 toUserId 和 sourceId 进行筛选，并按照创建时间倒序排列
	if err := common.DB.Where("(to_user_id = ? AND from_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", toUserId, sourceId, toUserId, sourceId).Order("created_at desc").First(&message).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return &message, nil
}

func ToRespMessage(message model.Message) model.RespMessage {
	var res model.RespMessage
	res.FromUserId = message.FromUserId
	res.ToUserId = message.ToUserId
	res.Content = message.Content
	res.Id = message.Id
	res.CreatedAt = message.CreatedAt.Unix()*1000 + 1000
	return res
}

func ToProtoMessage(msg model.Message) *pb.Message {
	return &pb.Message{
		Id:         msg.Id,
		ToUserId:   msg.ToUserId,
		Content:    msg.Content,
		CreateTime: strconv.Itoa(int(msg.CreatedAt.Unix())*1000 + 1000),
	}
}
