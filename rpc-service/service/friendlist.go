package service

import (
	"context"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
)

type FriendService struct {
	pb.UnimplementedFriendServiceServer
}

func (f FriendService) GetFriendList(ctx context.Context, req *pb.DouyinRelationFriendListRequest) (*pb.DouyinRelationFriendListResponse, error) {
	_, claim, _ := common.ParseToken(req.Token)
	sourceId := claim.UserId
	respUserList, err := service.GetFollowerList(sourceId, req.UserId)
	if err != nil {
		return &pb.DouyinRelationFriendListResponse{StatusCode: 1, StatusMsg: "请求好友列表失败！"}, nil
	}
	var friendList []*pb.FriendUser
	for _, u := range respUserList {
		message, err := service.GetLatestMessage(req.UserId, u.Id)
		if err != nil {
			return &pb.DouyinRelationFriendListResponse{StatusCode: 1, StatusMsg: "请求好友列表失败！"}, nil
		}
		var msgType int64
		if message.ToUserId == req.UserId {
			msgType = 0
		} else {
			msgType = 1
		}
		friendList = append(friendList, service.ToProtoFriend(u, message.Content, msgType))
	}
	return &pb.DouyinRelationFriendListResponse{StatusCode: 0, StatusMsg: "获取好友列表成功！", UserList: friendList}, nil
}
