package main

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/registry"
	"sync"
	"time"

	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
)

var wait sync.WaitGroup

func WatchService() {
	services := []string{"feed", "user", "publish", "favorite", "comment", "relation", "message", "friend"}

	for _, service := range services {
		wait.Add(1)
		go registry.WatchServiceName(service, &wait)
	}
}
func main() {
	//go service.RunMessageServer()

	WatchService()
	//需要先阻塞等watchServiceName初始化我们的douyinservice才可以初始化线程池
	wait.Wait()

	utils.InitConfig()
	common.InitDB()
	common.InitRedis()
	registry.InitAllConnPool()

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
