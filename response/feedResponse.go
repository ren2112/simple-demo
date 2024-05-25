package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FeedResponse struct {
	Response
	VideoList []*pb.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

func FeedResponseFun(c *gin.Context, response Response, respVideoList []*pb.Video, responseTime int64) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  response,
		VideoList: respVideoList,
		NextTime:  responseTime,
	})
}
