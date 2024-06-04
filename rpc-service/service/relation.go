package service

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"strconv"
	"sync"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

var LockRelationMap sync.Map

func (r RelationService) FollowAction(ctx context.Context, req *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)

	key := fmt.Sprintf("follow_to:%d", req.ToUserId)
	// 加载或创建锁，并确保最后解锁
	var mutex *sync.Mutex
	value, ok := LockRelationMap.Load(key)
	if !ok {
		mutex = new(sync.Mutex)
		LockRelationMap.Store(key, mutex)
	} else {
		mutex = value.(*sync.Mutex)
	}
	mutex.Lock()
	defer mutex.Unlock() // 确保在函数退出前解锁

	err := service.RelationAction(strconv.Itoa(int(req.ActionType)), claims.UserId, req.ToUserId)
	if err != nil {
		return &pb.DouyinRelationActionResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
	} else {
		return &pb.DouyinRelationActionResponse{StatusCode: 0, StatusMsg: "操作成功！"}, nil
	}
}

func (r RelationService) GetFollowList(ctx context.Context, req *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)
	//获取发起请求的用户id，为了判断这个用户是否有对别人的关注列表里面人的是否关注
	sourceId := claims.UserId

	respUserList, err := service.GetFollowList(sourceId, int(req.UserId))
	if err != nil {
		return &pb.DouyinRelationFollowListResponse{StatusCode: 1, StatusMsg: "获取关注列表失败！"}, nil
	} else {
		return &pb.DouyinRelationFollowListResponse{StatusCode: 0, StatusMsg: "获取关注列表成功！", UserList: respUserList}, nil
	}
}

func (r RelationService) GetFollowerList(ctx context.Context, req *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {
	_, claims, _ := common.ParseToken(req.Token)
	sourceId := claims.UserId
	respUserList, err := service.GetFollowerList(sourceId, req.UserId)
	if err != nil {
		return &pb.DouyinRelationFollowerListResponse{StatusCode: 1, StatusMsg: "获取粉丝列表失败！"}, nil
	} else {
		return &pb.DouyinRelationFollowerListResponse{StatusCode: 0, StatusMsg: "获取粉丝列表成功！", UserList: respUserList}, nil
	}
}
