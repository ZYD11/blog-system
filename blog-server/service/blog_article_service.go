package service

import (
	"blog-server/dao"
	"blog-server/models/blog"
)

type BlogArticleService struct {
	articleDao dao.BlogArticleDao
}

// 新增文章
func (s BlogArticleService) Add(article blog.Article) int64 {
	return s.articleDao.InsertArticle(article)
}

// 更新文章
func (s BlogArticleService) Update(article blog.Article) int64 {
	return s.articleDao.UpdateArticle(article)
}

// 删除文章
func (s BlogArticleService) Remove(articleId string) bool {
	return s.articleDao.DeleteArticle(articleId)
}

// 根据账号获取对应文章
func (s BlogArticleService) GetArticlesByUserId(userId int) *[]blog.ArticleInfo {
	return s.articleDao.GetArticlesByUserId(userId)
}

// 根据ID获取文章数据
func (s BlogArticleService) GetArticleById(articleId string) *blog.Article {
	return s.articleDao.GetArticleById(articleId)
}

// 根据多个ID获取文章数据
func (s BlogArticleService) GetArticlesByIds(articleId []string) *[]blog.ArticleInfo {
	return s.articleDao.GetArticlesByIds(articleId)
}

func (s BlogArticleService) Find(query blog.ArticleQuery) (*[]blog.Article, int64) {
	return s.articleDao.List(query)
}

func (s BlogArticleService) BathAudit(articleId string) int64 {
	return s.articleDao.UpdateAudit(articleId)
}
