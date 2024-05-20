package middleware

import (
	"github.com/RaymondCode/simple-demo/response"
	"github.com/gin-gonic/gin"
	"time"
)

func LeakBucketMiddleware(capacity int, rate time.Duration) gin.HandlerFunc {
	bucket := make(chan struct{}, capacity)
	go func() {
		ticker := time.Tick(rate)
		for {
			select {
			case <-ticker:
				select {
				case bucket <- struct{}{}:
				default:
				}
			}
		}
	}()
	return func(c *gin.Context) {
		select {
		case <-bucket:
			c.Next()
		default:
			response.CommonResp(c, 1, "请求过多！请稍后重试")
			c.Abort()
		}
	}
}
