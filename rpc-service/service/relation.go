package service

import (
	"context"
	"errors"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"strconv"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

func (r RelationService) FollowAction(ctx context.Context, req *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)
	err := service.RelationAction(strconv.Itoa(int(req.ActionType)), claims.UserId, int(req.ToUserId))
	if err != nil {
		return nil, errors.New("操作失败！")
	} else {
		return &pb.DouyinRelationActionResponse{StatusCode: 0, StatusMsg: "关注成功！"}, nil
	}
}

func (r RelationService) GetFollowList(ctx context.Context, req *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)
	//获取发起请求的用户id，为了判断这个用户是否有对别人的关注列表里面人的是否关注
	sourceId := claims.UserId

	respUserList, err := service.GetFollowList(sourceId, int(req.UserId))
	if err != nil {
		return nil, err
	} else {
		return &pb.DouyinRelationFollowListResponse{StatusCode: 0, StatusMsg: "获取关注列表成功！", UserList: respUserList}, nil
	}
}

func (r RelationService) GetFollowerList(ctx context.Context, req *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)
	sourceId := claims.UserId
	respUserList, err := service.GetFollowerList(sourceId, int(req.UserId))
	if err != nil {
		return nil, err
	} else {
		return &pb.DouyinRelationFollowerListResponse{StatusCode: 0, StatusMsg: "获取粉丝列表成功！", UserList: respUserList}, nil
	}
}
