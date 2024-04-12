package client

import (
	apiClient "blog-server/api/client"
	"github.com/gin-gonic/gin"
)

// 文章路由
func InitBlogCategoryRouter(router *gin.RouterGroup) {
	categoryApi := new(apiClient.CategoryApi)
	// 查询分类
	router.GET("/category", categoryApi.SearchCategory)         // 查询分类
	router.GET("/category/:id", categoryApi.SearchCategoryName) // 查询分类名
}
