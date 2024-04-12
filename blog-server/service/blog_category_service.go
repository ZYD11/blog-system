package service

import (
	"blog-server/dao"
	"blog-server/models/blog"
)

type BlogCategoryService struct {
	categoryDao dao.BlogCategoryDao
}

func (s BlogCategoryService) GetCategoryById(categoryId uint) *blog.Category {
	category := blog.Category{}
	category.Id = categoryId
	return s.categoryDao.GetCategoryCondition(category)
}

func (s BlogCategoryService) GetCategories() *[]blog.Category {
	return s.categoryDao.GetCategoryAll()
}
