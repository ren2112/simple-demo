package model

import "time"

type Favorite struct {
	Id         int64
	UserId     int64 `gorm:"foreignKey:User(id)"`
	User       User
	VideoId    int64 `gorm:"foreignKey:Video(id)"`
	Video      Video
	IsFavorite bool
	CreatedAt  time.Time
}
