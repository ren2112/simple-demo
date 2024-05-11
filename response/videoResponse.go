package response

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoListResponse struct {
	Response
	VideoList []model.RespVideo `json:"video_list"`
}

func VideoListResponseFun(c *gin.Context, response Response, resVideoList []model.RespVideo) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  response,
		VideoList: resVideoList,
	})
}
