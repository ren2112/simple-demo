package main

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"net"
	"time"
)

type VideoServer struct {
	pb.UnsafeVideoServiceServer
}

func (v VideoServer) GetFeedList(ctx context.Context, req *pb.DouyinFeedRequest) (*pb.DouyinFeedResponse, error) {
	_, claims, _ := common.ParseToken(*req.Token)

	//转化时间
	latestTime := *req.LatestTime
	var err error
	if latestTime == 0 {
		latestTime = time.Now().Unix() * 1000
	}
	if err != nil {
		return nil, err
	}
	var videoList = []*pb.Video{}
	latestTimeUTC := time.Unix(0, latestTime*int64(time.Millisecond))

	//获取视频流
	videoList, err = service.FeedVideoList(latestTimeUTC)
	if err != nil {
		return nil, err
	}

	//若用户登录了，需要判断视频作者是否关注以及是否对视频点赞
	if claims.UserId != 0 {
		for i, v := range videoList {
			//查找是否点赞
			var isFavorite bool
			isFavorite, err = service.JudgeFavorite(claims.UserId, *v.Id)
			if err != nil {
				return nil, err
			}
			*videoList[i].IsFavorite = isFavorite

			//	查找videoList的author里面is_follow
			// 查找作者是否被当前用户关注
			err = service.JudgeRelation(*v.Author.Id, claims.UserId)
			if err == nil {
				*videoList[i].Author.IsFollow = true
			} else if err == gorm.ErrRecordNotFound {
				*videoList[i].Author.IsFollow = false
			} else {
				return nil, err
			}
		}
	}
	responseTime := videoList[len(videoList)-1].CreatedAt
	return &pb.DouyinFeedResponse{StatusCode: proto.Int(0), StatusMsg: proto.String("获取成功！"), VideoList: videoList, NextTime: responseTime}, nil
}

func main() {
	common.InitDB()
	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	//在grpc服务端注册服务
	pb.RegisterVideoServiceServer(grpcServer, &VideoServer{})
	//	启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("启动grpc服务端失败！%v", err)
		return
	}
}
