package assist

import "github.com/RaymondCode/simple-demo/model"

func ToRespMessage(message model.Message) model.RespMessage {
	var res model.RespMessage
	res.FromUserId = message.FromUserId
	res.ToUserId = message.ToUserId
	res.Content = message.Content
	res.Id = message.Id
	res.CreatedAt = message.CreatedAt.Unix()*1000 + 1000
	return res
}
