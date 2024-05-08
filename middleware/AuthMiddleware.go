package middleware

import (
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Query("token")
		token, claims, err := common.ParseToken(tokenString)
		//token空或者不合法
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			})
			ctx.Abort()
			return
		}
		//校验user是否存在
		if user := assist.GetUserByID(int64(claims.UserId)); user.Id != 0 {
			ctx.Next()
		} else {
			ctx.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			})
			ctx.Abort()
		}

	}
}
