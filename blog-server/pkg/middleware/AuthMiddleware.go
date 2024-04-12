package middleware

import (
	"blog-server/dao"
	"blog-server/models/blog"
	"blog-server/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.Request.Header.Get("Authorization")

		// 非法token
		if tokenString == "" || len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		//if userData, _ := ctx.Get("user"); userData == nil {
		// 获取claims中的userId
		userId := claims.UserId
		userDao := dao.BlogUserDao{}
		user := userDao.GetUserByConditions(blog.User{Id: userId})
		// 将用户信息写入上下文便于读取
		ctx.Set("user", *user)
		// }
		ctx.Next()
	}
}
