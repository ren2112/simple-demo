package main

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/registry"
	"time"

	//"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
)

var signal = make(chan struct{})

func WatchService() {
	services := []string{"feed", "user", "publish", "favorite", "comment", "relation", "message", "friend"}

	for _, service := range services {
		go registry.WatchServiceName(service, signal)
	}

	for range services {
		<-signal
	}
}
func main() {
	//go service.RunMessageServer()

	WatchService()
	//需要先阻塞等watchServiceName初始化我们的douyinservice才可以初始化线程池

	utils.InitConfig()
	common.InitDB()
	common.InitRedis()
	common.InitAllConnPool()

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
