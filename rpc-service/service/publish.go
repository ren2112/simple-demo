package service

import (
	"context"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
)

type PublishService struct {
	pb.UnimplementedPublishServiceServer
}

func (p PublishService) GetPublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (*pb.DouyinPublishListResponse, error) {
	//获取发布列表
	var RespVideoList []*pb.Video
	RespVideoList, err := service.GetPublishVideoList(req.UserId)
	if err != nil {
		return &pb.DouyinPublishListResponse{StatusCode: 1, StatusMsg: "获取发布列表失败！"}, nil
	} else {
		return &pb.DouyinPublishListResponse{StatusCode: 0, StatusMsg: "获取发布列表成功！", VideoList: RespVideoList}, nil
	}
}
