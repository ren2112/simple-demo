package main

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/middleware"

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
	//使用漏桶来对请求节流
	r.Use(middleware.LeakBucketMiddleware(config.LEAKBUCKET_CAPACITY, config.LEAKBUCKET_RATE))
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
