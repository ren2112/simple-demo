package assist

import "github.com/RaymondCode/simple-demo/model"

func ToRespComment(comment model.Comment) model.RespComment { // 转换 Comment 中的 User 字段为 RespUser
	respComment := model.RespComment{
		Id:         comment.Id,
		User:       ToRespUser(comment.User),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化时间,
	}

	return respComment
}
