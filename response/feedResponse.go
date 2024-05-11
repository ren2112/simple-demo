package response

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FeedResponse struct {
	Response
	VideoList []model.RespVideo `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

func FeedResponseFun(c *gin.Context, respVideoList []model.RespVideo, responseTime int64) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: respVideoList,
		NextTime:  responseTime * 1000,
	})
}
