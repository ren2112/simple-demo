package assist

import "github.com/RaymondCode/simple-demo/model"

func ToRespVideo(video model.Video) model.RespVideo {
	respVideo := model.RespVideo{
		Id:            video.Id,
		Author:        ToRespUser(video.Author),
		Title:         video.Title,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
	}

	return respVideo
}
