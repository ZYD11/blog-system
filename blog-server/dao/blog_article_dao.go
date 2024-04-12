package dao

import (
	"blog-server/models/blog"
	"blog-server/pkg/page"
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
)

type BlogArticleDao struct {
}

func (d BlogArticleDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table(blog.Article{}.TableName())
}

// 列表数据查询
func (d BlogArticleDao) List(query blog.ArticleQuery) (*[]blog.Article, int64) {
	articles := make([]blog.Article, 0)
	session := d.sql(SqlDB.NewSession())

	if gotool.StrUtils.HasNotEmpty(query.Keyword) {
		session.And("(title LIKE concat('%', ?, '%') OR content LIKE concat('%', ?, '%'))", query.Keyword, query.Keyword)
	}

	if gotool.StrUtils.HasNotEmpty(query.CategoryId) && query.CategoryId != "0" {
		session.And("category_id=?", query.CategoryId)
	}

	if gotool.StrUtils.HasNotEmpty(query.BeginTime) {
		session.And("date_format(created_at,'%y%m%d') >= date_format(?,'%y%m%d')", query.BeginTime)
	}
	if gotool.StrUtils.HasNotEmpty(query.EndTime) {
		session.And("date_format(created_at,'%y%m%d') <= date_format(?,'%y%m%d')", query.EndTime)
	}

	total, _ := page.GetTotal(session.Clone())

	session.Select("id,category_id,title,LEFT(content,80) AS content,head_image,created_at,user_id")
	err := session.Limit(query.PageSize, page.StartSize(query.PageNum, query.PageSize)).OrderBy("created_at DESC").Find(&articles)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil, 0
	}

	return &articles, total
}

// InsertArticle 添加文章
func (d BlogArticleDao) InsertArticle(article blog.Article) int64 {
	session := SqlDB.NewSession()
	session.Begin()

	insert, err := session.Insert(&article)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()

	return insert
}

// UpdateArticle 修改文章
func (d BlogArticleDao) UpdateArticle(article blog.Article) int64 {
	session := SqlDB.NewSession()
	session.Begin()
	update, err := session.Where("id = ?", article.Id).Update(&article)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()
	return update
}

// 删除文章
func (d BlogArticleDao) DeleteArticle(articleId string) bool {
	session := SqlDB.NewSession()
	session.Begin()

	ids := strings.Split(articleId, ",")
	list := make([]int64, 0)
	for _, id := range ids {
		parseInt, _ := strconv.ParseInt(id, 10, 64)
		list = append(list, parseInt)
	}

	_, err := session.In("id", list).Delete(&blog.Article{})
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return false
	}
	session.Commit()
	return true
}

// 根据用条件查询用户数据
func (d BlogArticleDao) GetArticlesByUserId(userId int) *[]blog.ArticleInfo {
	var articles []blog.ArticleInfo
	session := d.sql(SqlDB.NewSession())
	session.Select("id,category_id,title,LEFT(content,80) AS content, head_image, created_at")
	err := session.Where("user_id = ?", userId).OrderBy("created_at desc").Find(&articles)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &articles
}

// 根据用条件查询用户数据
func (d BlogArticleDao) GetArticlesByIds(articleId []string) *[]blog.ArticleInfo {
	var articles []blog.ArticleInfo
	session := d.sql(SqlDB.NewSession())
	session.Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at")
	err := session.In("id", articleId).OrderBy("created_at desc").Find(&articles)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &articles
}

func (d BlogArticleDao) GetArticleById(articleId string) *blog.Article {
	var article blog.Article
	session := d.sql(SqlDB.NewSession())
	_, err := session.Where("id = ?", articleId).Get(&article)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &article
}

func (d BlogArticleDao) GetArticleByConditions(article blog.Article) *blog.Article {
	i, err := SqlDB.Get(&article)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	if i {
		return &article
	}
	return nil
}

func (d BlogArticleDao) UpdateAudit(articleId string) int64 {
	ids := strings.Split(articleId, ",")
	list := make([]int64, 0)
	for _, id := range ids {
		parseInt, _ := strconv.ParseInt(id, 10, 64)
		list = append(list, parseInt)
	}

	session := SqlDB.NewSession()
	session.Begin()
	update, err := session.In("id", ids).Update(&blog.Article{Status: 2})

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()
	return update
}
