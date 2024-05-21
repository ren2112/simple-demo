package model

type User struct {
	Id              int64  `json:"id,omitempty"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
	Name            string `json:"name,omitempty" gorm:"index;unique"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"-"`
	Password        string `json:"password" gorm:"size:255,not null"`
	WorkCount       int    `json:"work_count" gorm:"default:0"`
	TotalFavorited  int    `json:"total_favorited" gorm:"default:0"`
	FavoriteCount   int    `json:"favorite_count" gorm:"default:0"`
}

type RespUser struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	WorkCount       int    `json:"work_count"`
	TotalFavorited  int    `json:"total_favorited"`
	FavoriteCount   int    `json:"favorite_count"`
}
