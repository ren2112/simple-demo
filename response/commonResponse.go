package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func CommonResp(c *gin.Context, statusCode int32, statusMsg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

func CommonServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		StatusCode: 1,
		StatusMsg:  "操作失败！",
	})
}

func RPCServerUnstart(c *gin.Context, str string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: 1,
		StatusMsg:  str + "服务未启动",
	})
}
