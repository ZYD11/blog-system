package dao

import (
	"blog-server/models/blog"
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
)

type BlogCategoryDao struct {
}

func (d BlogCategoryDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table(blog.Category{}.TableName())
}

func (d BlogCategoryDao) GetCategoryCondition(category blog.Category) *blog.Category {
	i, err := SqlDB.Get(&category)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	if i {
		return &category
	}
	return nil
}

func (d BlogCategoryDao) GetCategoryAll() *[]blog.Category {
	var categories []blog.Category
	session := d.sql(SqlDB.NewSession())
	err := session.Find(&categories)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &categories
}
