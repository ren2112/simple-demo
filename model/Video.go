package model

import "time"

type Video struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt time.Time
	AuthorID  int64 `json:"author_id,omitempty" gorm:"column:author_id;index;references:User"`
	Author    User  `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	//点赞的情况时候需要用到多对多
	Users         []User `gorm:"many2many:user_videos;association_autoupdate:false"`
	Title         string `json:"title"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count" gorm:"default:0"`
	CommentCount  int64  `json:"comment_count" gorm:"default:0"`
	IsFavorite    bool   `json:"is_favorite" gorm:"-"`
}

type RespVideo struct {
	Id            int64    `json:"id,omitempty"`
	Author        RespUser `json:"author,omitempty"`
	Title         string   `json:"title"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
}
