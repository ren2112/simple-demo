package response

import (
	pb "github.com/RaymondCode/simple-demo/controller/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FeedResponse struct {
	Response
	VideoList []*pb.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

func FeedResponseFun(c *gin.Context, respVideoList []*pb.Video, responseTime int64) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: respVideoList,
		NextTime:  responseTime * 1000,
	})
}
