package assist

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
)

func IsFollowed(aUserId, bUserId int64) bool {
	// 查询是否a关注了指定b用户的逻辑
	follow := model.Follow{}
	common.DB.Where("user_id=? and follower_user_id=?", bUserId, aUserId).Find(&follow)
	if follow.Id != 0 {
		return true
	}
	return false
}
