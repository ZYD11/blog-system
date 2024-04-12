package client

import (
	apiClient "blog-server/api/client"
	"blog-server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// 用户路由
func InitBlogUserRouter(router *gin.RouterGroup) {
	userApi := new(apiClient.UserApi)
	fileApi := new(apiClient.FileApi)
	// 登录
	router.POST("/login", userApi.Login)
	// 注册
	router.POST("/register", userApi.Register)

	router.POST("/logout", userApi.LogOut)

	// 上传图像
	router.POST("/upload", fileApi.Upload)
	router.POST("/upload/rich_editor_upload", fileApi.RichEditorUpload)

	userRoutes := router.Group("/user")
	{
		userRoutes.Use(middleware.AuthMiddleware())
		userRoutes.GET("userInfo", userApi.GetUserInfo)             // 验证用户
		userRoutes.GET("briefInfo/:id", userApi.GetBriefInfo)       // 获取用户简要信息
		userRoutes.GET("detailedInfo/:id", userApi.GetDetailedInfo) // 获取用户详细信息
		userRoutes.PUT("avatar/:id", userApi.ModifyAvatar)          // 修改头像
		userRoutes.PUT("name/:id", userApi.ModifyName)              // 修改用户名
	}

	// 我的收藏
	colRoutes := router.Group("/collects")
	{
		colRoutes.Use(middleware.AuthMiddleware())
		colRoutes.GET(":id", userApi.Collects)        // 查询收藏
		colRoutes.PUT("new/:id", userApi.NewCollect)  // 收藏
		colRoutes.DELETE(":index", userApi.UnCollect) // 取消收藏
	}
	// 我的关注
	folRoutes := router.Group("/following")
	{
		folRoutes.Use(middleware.AuthMiddleware())
		folRoutes.GET(":id", userApi.Following)      // 查询关注
		folRoutes.PUT("new/:id", userApi.NewFollow)  // 关注
		folRoutes.DELETE(":index", userApi.UnFollow) // 取消关注
	}
}
