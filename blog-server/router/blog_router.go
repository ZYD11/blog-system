package router

import (
	apiAdmin "blog-server/api/admin"
	"github.com/gin-gonic/gin"
)

func initBlogRouter(router *gin.RouterGroup) {
	articleApi := new(apiAdmin.ArticleApi)
	articleGroup := router.Group("/article")
	{
		articleGroup.GET("/list", articleApi.List)
		articleGroup.GET("/:id", articleApi.Get)
		articleGroup.DELETE("/:ids", articleApi.Delete)
		articleGroup.POST("/audit", articleApi.Audit)
	}
}
