package controller

import (
	"github.com/RaymondCode/simple-demo/response"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	//连接grpc服务端
	conn, err := grpc.Dial("127.0.0.1:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接：%v", err)
	}
	defer conn.Close()

	//	建立连接
	client := pb.NewVideoFeedServiceClient(conn)
	resp, err := client.GetFeedList(c, &pb.DouyinFeedRequest{LatestTime: latestTime, Token: tokenStr})
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	response.FeedResponseFun(c, resp.VideoList, resp.NextTime)
}
