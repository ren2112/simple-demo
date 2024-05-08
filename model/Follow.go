package model

import "time"

type Follow struct {
	Id             int64
	UserId         int64 `gorm:"foreignKey:User(id)"`
	User           User
	FollowerUserId int64 `gorm:"foreignKey:User(id)"`
	FollowerUser   User
	CreatedAt      time.Time
}
