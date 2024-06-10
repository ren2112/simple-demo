package middleware

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Query("token")
		if tokenString == "" {
			tokenString = ctx.PostForm("token")
		}
		token, claims, err := common.ParseToken(tokenString)
		//token空或者不合法
		if err != nil || !token.Valid {
			response.CommonResp(ctx, 1, "请登录或重新登陆")
			ctx.Abort()
			return
		}

		//如果上下文存在用户就返回
		if _, ok := ctx.Get("user"); ok {
			ctx.Next()
			return
		}

		//检查token是否在redis的黑名单里面
		if common.CheckTokenInBlacklist(ctx, tokenString) {
			response.CommonResp(ctx, 1, "你被加入黑名单！24h后解封！")
		}

		//尝试从缓存获得用户信息
		cachedUser, err := common.GetCachedUser(ctx, claims.UserName)
		if err == nil && cachedUser != nil {
			ctx.Set("user", *cachedUser)
			ctx.Next()
			return
		}

		//缓存没有数据
		if user := service.GetUserByID(claims.UserId); user.Id != 0 {
			err = common.CacheUser(ctx, claims.UserName, user)
			if err != nil {
				fmt.Println("缓存用户失败！", err)
			}
			ctx.Set("user", user)
			ctx.Next()
		} else {
			response.CommonResp(ctx, 1, "请登录或重新登陆")
			ctx.Abort()
		}
	}
}
