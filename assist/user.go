package assist

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
)

func GetUserByName(userName string) (user model.User) {
	common.DB.Where("name=?", userName).First(&user)
	return user
}

func GetUserByID(ID int64) (user model.User) {
	common.DB.Find(&user, ID)
	return user
}
