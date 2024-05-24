package service

import (
	"context"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"gorm.io/gorm"
	"time"
)

type VideoFeedService struct {
	pb.UnimplementedVideoFeedServiceServer
}

func (v VideoFeedService) GetFeedList(ctx context.Context, req *pb.DouyinFeedRequest) (*pb.DouyinFeedResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)

	//转化时间
	latestTime := req.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}

	var videoList = []*pb.Video{}

	//获取视频流
	videoList, err := service.FeedVideoList(latestTime)
	if err != nil {
		return nil, err
	}

	//若用户登录了，需要判断视频作者是否关注以及是否对视频点赞
	if claims.UserId != 0 {
		for i, v := range videoList {
			//查找是否点赞
			var isFavorite bool
			isFavorite, err = service.JudgeFavorite(claims.UserId, v.Id)
			if err != nil {
				return nil, err
			}
			videoList[i].IsFavorite = isFavorite

			//	查找videoList的author里面is_follow
			// 查找作者是否被当前用户关注
			err = service.JudgeRelation(v.Author.Id, claims.UserId)
			if err == nil {
				videoList[i].Author.IsFollow = true
			} else if err == gorm.ErrRecordNotFound {
				videoList[i].Author.IsFollow = false
			} else {
				return nil, err
			}
		}
	}
	var responseTime int64
	if len(videoList) == 0 {
		responseTime = time.Now().Unix()
	} else {
		responseTime = videoList[len(videoList)-1].CreatedAt
	}
	return &pb.DouyinFeedResponse{StatusCode: 0, StatusMsg: "获取成功！", VideoList: videoList, NextTime: responseTime}, nil
}
