package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"strconv"
	"time"
)

func CreateMessage(toUserId string, content string, userId int64) error {
	userIdTarget, _ := strconv.Atoi(toUserId)
	var message model.Message

	message.Content = content
	message.FromUserId = userId
	message.ToUserId = int64(userIdTarget)
	if err := common.DB.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

func GetMessageList(userIdTarget string, fromUserId int64, preMsgTimeUTC time.Time) ([]model.RespMessage, error) {
	var messageList []model.Message
	if err := common.DB.Where("(to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)", userIdTarget, fromUserId, fromUserId, userIdTarget).
		Where("created_at > ?", preMsgTimeUTC).
		Order("created_at").
		Find(&messageList).Error; err != nil {
		return nil, err
	}

	var resMessageList []model.RespMessage
	for _, v := range messageList {
		var resMessage = ToRespMessage(v)
		resMessageList = append(resMessageList, resMessage)
	}
	return resMessageList, nil
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
