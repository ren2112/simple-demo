package middleware

import (
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
		//校验user是否存在
		if user := service.GetUserByID(claims.UserId); user.Id != 0 {
			ctx.Set("user", user)
			ctx.Next()
		} else {
			response.CommonResp(ctx, 1, "请登录或重新登陆")
			ctx.Abort()
		}
	}
}
