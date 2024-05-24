package main

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/middleware"
	"time"

	//"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	//go service.RunMessageServer()

	utils.InitConfig()
	common.InitDB()
	common.InitRedis()

	r := gin.Default()
	//使用token桶来对请求节流
	bucket := make(chan struct{}, config.TOKENBUCKET_CAPACITY)
	go func() {
		ticker := time.Tick(config.TOKENBUCKET_RATE)
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
	r.Use(middleware.TokenBucketMiddleware(&bucket))
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
