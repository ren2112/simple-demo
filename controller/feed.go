package controller

import (
	"github.com/RaymondCode/simple-demo/registry"
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	if latestTimeStr == "" {
		latestTimeStr = "0"
	}
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		response.CommonResp(c, 1, "请求时间错误！")
		return
	}
	tokenStr := c.Query("token")

	// 从连接池中获取连接
	connPool, ok := registry.GetPool("feed")
	if !ok {
		response.RPCServerUnstart(c, "feed")
		return
	}
	conn := connPool.Get()

	// 建立连接
	client := pb.NewVideoFeedServiceClient(conn)
	resp, err := client.GetFeedList(c, &pb.DouyinFeedRequest{LatestTime: latestTime, Token: tokenStr})
	connPool.Put(conn)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.FeedResponseFun(c, response.Response{StatusCode: resp.StatusCode, StatusMsg: resp.StatusMsg}, resp.VideoList, resp.NextTime)
}
