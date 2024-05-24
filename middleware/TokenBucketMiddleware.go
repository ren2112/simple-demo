package middleware

import (
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
)

func TokenBucketMiddleware(bucket *chan struct{}) gin.HandlerFunc {

	return func(c *gin.Context) {
		select {
		case <-*bucket:
			c.Next()
		default:
			response.CommonResp(c, 1, "请求过多！请稍后重试")
			c.Abort()
		}
	}
}
