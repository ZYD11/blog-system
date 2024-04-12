package client

import (
	"blog-server/models/blog"
	"blog-server/service"
	response "blog-server/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CategoryApi 博客分类操作api
type CategoryApi struct {
	categoryService service.BlogCategoryService
}

func (categoryApi CategoryApi) SearchCategory(ctx *gin.Context) {
	var categories *[]blog.Category
	categories = categoryApi.categoryService.GetCategories()

	if categories == nil {
		response.Fail(ctx, nil, "查找失败")
		return
	}
	response.Success(ctx, gin.H{"categories": categories}, "查找成功")
}

func (categoryApi CategoryApi) SearchCategoryName(ctx *gin.Context) {
	var category *blog.Category
	categoryId := ctx.Params.ByName("id")
	paramCategory, _ := strconv.Atoi(categoryId)
	category = categoryApi.categoryService.GetCategoryById(uint(paramCategory))

	if category == nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	response.Success(ctx, gin.H{"categoryName": category.CategoryName}, "查找成功")
}
