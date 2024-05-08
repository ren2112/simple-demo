package model

import "time"

type Favorite struct {
	Id        int64
	UserId    int64 `gorm:"foreignKey:User(id)"`
	VideoId   int64 `gorm:"foreignKey:Video(id)"`
	Video     Video
	CreatedAt time.Time
}
