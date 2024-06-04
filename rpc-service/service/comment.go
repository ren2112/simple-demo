package service

import (
	"context"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
}

func (c CommentService) CommentAction(ctx context.Context, req *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	respComment := pb.Comment{}
	err := service.CommentAction(req.ActionType, req.User, req.VideoId, req.CommentText, req.CommentId, &respComment)
	if err != nil {
		return &pb.DouyinCommentActionResponse{StatusCode: 1, StatusMsg: err.Error()}, nil
	}
	return &pb.DouyinCommentActionResponse{StatusCode: 0, StatusMsg: "评论成功！", Comment: &respComment}, nil
}

func (c CommentService) GetCommentList(ctx context.Context, req *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	respComments, err := service.GetCommentList(req.VideoId)
	if err != nil {
		return &pb.DouyinCommentListResponse{StatusCode: 1, StatusMsg: "获取评论列表失败！"}, nil
	}
	return &pb.DouyinCommentListResponse{StatusCode: 0, StatusMsg: "获取评论成功！", CommentList: respComments}, nil
}
