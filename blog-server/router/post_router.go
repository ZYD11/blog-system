package router

import (
	apiAdmin "blog-server/api/admin"
	"github.com/gin-gonic/gin"
)

// 初始化岗位路由
func initPostRouter(router *gin.RouterGroup) {
	v := new(apiAdmin.PostApi)
	group := router.Group("/post")
	{
		//查询岗位数据
		group.GET("/list", v.List)
		//添加岗位
		group.POST("/add", v.Add)
		//查询岗位详情
		group.GET("/:postId", v.Get)
		//删除岗位数据
		group.DELETE("/:postId", v.Delete)
		//修改岗位数据接口
		group.PUT("/edit", v.Edit)
		//导出excel
		group.GET("/export", v.Export)
	}
}
