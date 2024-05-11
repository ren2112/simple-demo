package response

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListResponse struct {
	Response
	CommentList []model.RespComment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment model.RespComment `json:"comment,omitempty"`
}

func CommentListResponseFun(c *gin.Context, response Response, respComments []model.RespComment) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: respComments,
	})
}

func CommentActionResponseFun(c *gin.Context, response Response, respComment model.RespComment) {
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: response,
		Comment:  respComment,
	})
}
