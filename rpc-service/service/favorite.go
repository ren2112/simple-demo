package service

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"sync"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
}

var LockFavoriteMap sync.Map

func (f FavoriteService) FavoriteAction(ctx context.Context, req *pb.DouyinFavoriteActionRequest) (*pb.DouyinFavoriteActionResponse, error) {
	_, claim, _ := common.ParseToken(req.Token)
	userId := claim.UserId
	key := fmt.Sprintf("favorite_video:%d", req.VideoId)

	// 加载或创建锁，并确保最后解锁
	var mutex *sync.Mutex
	value, ok := LockFavoriteMap.Load(key)
	if !ok {
		mutex = new(sync.Mutex)
		LockFavoriteMap.Store(key, mutex)
	} else {
		mutex = value.(*sync.Mutex)
	}
	mutex.Lock()
	defer mutex.Unlock() // 确保在函数退出前解锁

	err := service.FavoriteAction(req.ActionType, userId, req.VideoId)
	if err != nil {
		return &pb.DouyinFavoriteActionResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
	}
	return &pb.DouyinFavoriteActionResponse{StatusCode: 0, StatusMsg: "操作成功！"}, nil
}
func (f FavoriteService) GetFavoriteList(ctx context.Context, req *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	_, claim, _ := common.ParseToken(req.Token)
	sourceId := claim.UserId
	resVideoList, err := service.GetFavoriteList(sourceId, req.UserId)
	if err != nil {
		return &pb.DouyinFavoriteListResponse{StatusCode: 1, StatusMsg: "获取点赞列表异常！"}, nil
	}
	return &pb.DouyinFavoriteListResponse{StatusCode: 0, StatusMsg: "获取成功！", VideoList: resVideoList}, nil
}
