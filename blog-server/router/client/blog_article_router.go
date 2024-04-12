package client

import (
	apiClient "blog-server/api/client"
	"blog-server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// 文章路由
func InitBlogArticleRouter(router *gin.RouterGroup) {
	articleApi := new(apiClient.ArticleApi)
	//用户文章的增删查改
	articleRoutes := router.Group("/article")
	articleRoutes.Use(middleware.AuthMiddleware())
	{
		articleRoutes.POST("", articleApi.Create)                       // 发布文章
		articleRoutes.PUT(":id", articleApi.Update)                     // 修改文章
		articleRoutes.DELETE(":id", articleApi.Delete)                  // 删除文章
		articleRoutes.GET(":id", articleApi.Show)                       // 查看文章
		articleRoutes.POST("/list", articleApi.List)                    // 显示文章列表
		articleRoutes.POST("/comment_list", articleApi.CommentList)     // 显示评论列表
		articleRoutes.POST("/addComment", articleApi.NewComment)        // 添加新评论
		articleRoutes.POST("/delComment/:id", articleApi.DeleteComment) // 删除新评论
		articleRoutes.GET("/like/:id", articleApi.Like)                 // 是否点赞
		articleRoutes.PUT("/newLike/:id", articleApi.NewLike)           // 新增点赞
		articleRoutes.DELETE("/unLike/:id", articleApi.UnLike)          // 取消点赞
	}
}
