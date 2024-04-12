package router

import (
	"blog-server/pkg/jwt"
	"blog-server/pkg/middleware"
	"blog-server/pkg/middleware/logger"
	. "blog-server/router/client"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) *gin.Engine {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(logger.LoggerToFile())
	router.Use(middleware.Recover)
	router.Use(middleware.CORSMiddleware())

	// router.Use(filter.DemoHandler())
	//后台管理api
	adminRouter := router.Group("/api/admin")
	adminRouter.Use(jwt.JWTAuth())
	{
		//登录接口
		initLoginRouter(adminRouter)
		//用户路由接口
		initUserRouter(adminRouter)
		//部门路由注册
		initDeptRouter(adminRouter)
		//初始化字典数据路由
		initDictDataRouter(adminRouter)
		//注册配置路由
		initConfigRouter(adminRouter)
		//注册角色路由
		initRoleRouter(adminRouter)
		//注册菜单路由
		initMenuRouter(adminRouter)
		//注册岗位路由
		initPostRouter(adminRouter)
		//注册字典类型路由
		initDictTypeRouter(adminRouter)
		//注册公告路由
		initNoticeRouter(adminRouter)
		
		initBlogRouter(adminRouter)
	}

	clientRouter := router.Group("/")
	{
		InitBlogUserRouter(clientRouter)
		InitBlogArticleRouter(clientRouter)
		InitBlogCategoryRouter(clientRouter)
	}
	return router
}
