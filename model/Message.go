package model

import "time"

type Message struct {
	Id         int64     `json:"id"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
}

type RespMessage struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"create_time"`
}
