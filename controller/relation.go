package controller

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	targetID, err := strconv.Atoi(c.Query("to_user_id"))
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		response.CommonResp(c, 1, "操作失败！")
		return
	}

	conn := common.GetRelationConnection()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.FollowAction(c, &pb.DouyinRelationActionRequest{Token: token, ActionType: int32(actionType), ToUserId: int64(targetID)})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.CommonResp(c, resp.StatusCode, resp.StatusMsg)

	//if err != nil {
	//	response.CommonResp(c, 1, "用户不存在")
	//	return
	//}
	//err = service.RelationAction(actionType, user.(model.User).Id, targetID)
	//if err != nil {
	//	response.CommonResp(c, 1, err.Error())
	//} else {
	//	response.CommonResp(c, 0, "操作成功")
	//}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败"}, nil)
		return
	}

	conn := common.GetRelationConnection()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.GetFollowList(c, &pb.DouyinRelationFollowListRequest{UserId: int64(userId), Token: token})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserListResponseFun(c, response.Response{StatusCode: 0}, resp.UserList)

	//获取发起请求的用户id，为了判断这个用户是否有对别人的关注列表里面人的是否关注
	//sourceId := claims.UserId
	//if sourceId == 0 {
	//	response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	//}
	//
	//respUserList, err := service.GetFollowList(sourceId, userId)
	//if err != nil {
	//	response.CommonResp(c, 1, err.Error())
	//} else {
	//	response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
	//}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
		return
	}

	conn := common.GetRelationConnection()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.GetFollowerList(c, &pb.DouyinRelationFollowerListRequest{UserId: int64(userId), Token: token})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserListResponseFun(c, response.Response{StatusCode: 0}, resp.UserList)

	////获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	//sourceId := claims.UserId
	//if sourceId == 0 {
	//	response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	//}
	//
	//respUserList, err := service.GetFollowerList(sourceId, userId)
	//if err != nil {
	//	response.CommonResp(c, 1, err.Error())
	//} else {
	//	response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
	//}
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	token := c.Query("token")
	_, claims, _ := common.ParseToken(token)
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
		return
	}

	//获取发起请求的用户id，为了判断这个用户是否有对别人的粉丝列表里面的粉丝是否关注
	sourceId := claims.UserId
	if sourceId == 0 {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "你还没登陆哦！"}, nil)
	}

	respUserList, err := service.GetFollowerList(sourceId, userId)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	} else {
		response.UserListResponseFun(c, response.Response{StatusCode: 0}, respUserList)
	}
}
