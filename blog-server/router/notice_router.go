package router

import (
	apiAdmin "blog-server/api/admin"
	"github.com/gin-gonic/gin"
)

func initNoticeRouter(router *gin.RouterGroup) {
	v := new(apiAdmin.NoticeApi)
	group := router.Group("/notice")
	{
		group.GET("/list", v.List)
		//添加公告
		group.POST("/add", v.Add)
		//删除
		group.DELETE("/:ids", v.Delete)
		//查询
		group.GET("/:id", v.Get)
		//修改
		group.PUT("/edit", v.Edit)
	}
}
