package model

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"time"
)

type Comment struct {
	Id        int64
	VideoId   int64 `gorm:"foreignKey:Video(id)"`
	User      User
	UserId    int64 `gorm:"foreignKey:User(id)"`
	Content   string
	CreatedAt time.Time
}

type RespComment struct {
	Id         int64    `json:"id"`
	User       *pb.User `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}
