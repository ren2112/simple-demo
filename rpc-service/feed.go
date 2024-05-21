package main

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
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

func main() {
	utils.InitConfig()
	common.InitDB()
	listen, err := net.Listen("tcp", ":9091")
	if err != nil {
		fmt.Printf("无法启动监听：%v\n", err)
		return
	}

	// 创建 gRPC 服务器对象
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	// 在 gRPC 服务器上注册服务
	pb.RegisterVideoServiceServer(grpcServer, &VideoServer{})

	// 启动 gRPC 服务
	fmt.Println("启动 gRPC 服务...")
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("启动 gRPC 服务失败：%v\n", err)
		return
	}
}
