package model

import "time"

type User struct {
	Id              int64     `json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name,omitempty"`
	Avatar          string    `json:"avatar"`
	BackgroundImage string    `json:"background_image"`
	Signature       string    `json:"signature"`
	FollowCount     int64     `json:"follow_count,omitempty"`
	FollowerCount   int64     `json:"follower_count,omitempty"`
	IsFollow        bool      `json:"is_follow,omitempty" gorm:"-"`
	Password        string    `json:"password" gorm:"size:255,not null"`
	Salt            string    `json:"salt"`
	Videos          []Video   `gorm:"many2many:user_videos;association_autoupdate:false"`
	WorkCount       int       `json:"work_count" gorm:"default:0"`
	TotalFavorited  int       `json:"total_favorited" gorm:"default:0"`
	FavoriteCount   int       `json:"favorite_count" gorm:"default:0"`
}
