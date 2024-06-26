package router

import (
	apiAdmin "blog-server/api/admin"
	"github.com/gin-gonic/gin"
)

func initConfigRouter(router *gin.RouterGroup) {
	v := new(apiAdmin.ConfigApi)
	group := router.Group("/config")
	{
		//根据参数键名查询参数值
		group.GET("/configKey/:configKey", v.GetConfigKey)
		//查询设置列表
		group.GET("/list", v.List)
		//添加
		group.POST("/add", v.Add)
		//查询
		group.GET("/:configId", v.Get)
		//修改
		group.PUT("/edit", v.Edit)
		//批量删除
		group.DELETE("/:ids", v.Delete)
		//刷新缓存
		group.DELETE("/refreshCache", v.RefreshCache)
		//导出数据
		group.GET("/export", v.Export)
	}
}
