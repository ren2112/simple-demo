package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoListResponse struct {
	Response
	VideoList []*pb.Video `json:"video_list"`
}

func VideoListResponseFun(c *gin.Context, response Response, resVideoList []*pb.Video) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  response,
		VideoList: resVideoList,
	})
}
