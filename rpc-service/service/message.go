package service

import (
	"context"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"time"
)

type MessageService struct {
	pb.UnimplementedChatServiceServer
}

func (m MessageService) ChatAction(ctx context.Context, req *pb.DouyinChatActionRequest) (*pb.DouyinChatActionResponse, error) {
	if req.Content == "" {
		return &pb.DouyinChatActionResponse{StatusCode: 1, StatusMsg: "发送消息不能为空"}, nil
	}
	_, claim, _ := common.ParseToken(req.Token)

	//添加信息于数据库
	err := service.CreateMessage(req.ToUserId, req.Content, claim.UserId)
	if err != nil {
		return &pb.DouyinChatActionResponse{StatusCode: 1, StatusMsg: "发送消息失败"}, nil
	}
	return &pb.DouyinChatActionResponse{StatusCode: 0, StatusMsg: "发送消息成功！"}, nil
}

func (m MessageService) GetChatList(ctx context.Context, req *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	_, claim, _ := common.ParseToken(req.Token)
	fromUserId := claim.UserId

	preMsgTimeUTC := time.Unix(0, req.PreMsgTime*int64(time.Millisecond))
	var resMessageList []*pb.Message
	resMessageList, err := service.GetMessageList(req.ToUserId, fromUserId, preMsgTimeUTC)
	if err != nil {
		return &pb.DouyinMessageChatResponse{StatusCode: 1, StatusMsg: "获取聊天记录失败！"}, nil
	} else {
		return &pb.DouyinMessageChatResponse{StatusCode: 0, StatusMsg: "获取聊天记录成功!", MessageList: resMessageList}, nil
	}
}
