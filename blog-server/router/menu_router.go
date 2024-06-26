package router

import (
	apiAdmin "blog-server/api/admin"
	"github.com/gin-gonic/gin"
)

func initMenuRouter(router *gin.RouterGroup) {
	v := new(apiAdmin.MenuApi)
	vg := router.Group("/menu")
	{
		//查询菜单数据
		vg.GET("/:menuId", v.GetInfo)
		//查询菜单列表
		vg.GET("/list", v.List)
		//加载对应角色菜单列表树
		vg.GET("/roleMenuTreeselect/:roleId", v.RoleMenuTreeSelect)
		//获取菜单下拉树列表
		vg.GET("treeselect", v.TreeSelect)
		//添加菜单数据
		vg.POST("/add", v.Add)
		//修改菜单数据
		vg.PUT("/edit", v.Edit)
		//删除菜单数据
		vg.DELETE("/:menuId", v.Delete)
	}
}
