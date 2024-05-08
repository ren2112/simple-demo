package model

type UserVideo struct {
	UserID     int64 `json:"user_id"`
	VideoID    int64 `json:"video_id"`
	IsFavorite bool  `json:"is_favorite"`
}
