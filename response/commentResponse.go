package response

import (
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListResponse struct {
	Response
	CommentList []*pb.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment *pb.Comment `json:"comment,omitempty"`
}

func CommentListResponseFun(c *gin.Context, response Response, respComments []*pb.Comment) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    response,
		CommentList: respComments,
	})
}

func CommentActionResponseFun(c *gin.Context, response Response, respComment *pb.Comment) {
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: response,
		Comment:  respComment,
	})
}
