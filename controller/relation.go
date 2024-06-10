package controller

import (
	"github.com/RaymondCode/simple-demo/registry"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
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

	connPool, ok := registry.GetPool("relation")
	if !ok {
		response.RPCServerUnstart(c, "relation")
		return
	}
	conn := connPool.Get()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.FollowAction(c, &pb.DouyinRelationActionRequest{Token: token, ActionType: int32(actionType), ToUserId: int64(targetID)})
	connPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.CommonResp(c, resp.StatusCode, resp.StatusMsg)
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败"}, nil)
		return
	}

	connPool, ok := registry.GetPool("relation")
	if !ok {
		response.RPCServerUnstart(c, "relation")
		return
	}
	conn := connPool.Get()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.GetFollowList(c, &pb.DouyinRelationFollowListRequest{UserId: int64(userId), Token: token})
	connPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserListResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.UserList)
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
		return
	}

	connPool, ok := registry.GetPool("relation")
	if !ok {
		response.RPCServerUnstart(c, "relation")
		return
	}
	conn := connPool.Get()

	client := pb.NewRelationServiceClient(conn)
	resp, err := client.GetFollowerList(c, &pb.DouyinRelationFollowerListRequest{UserId: int64(userId), Token: token})
	connPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.UserListResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.UserList)
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	token := c.Query("token")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.UserListResponseFun(c, response.Response{StatusCode: 1, StatusMsg: "操作失败！"}, nil)
		return
	}

	connPool, ok := registry.GetPool("friend")
	if !ok {
		response.RPCServerUnstart(c, "friend")
		return
	}
	conn := connPool.Get()

	client := pb.NewFriendServiceClient(conn)
	resp, err := client.GetFriendList(c, &pb.DouyinRelationFriendListRequest{UserId: int64(userId), Token: token})
	connPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.FriendListResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.UserList)
}
