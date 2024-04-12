package dao

import (
	"blog-server/models/blog"
	"blog-server/pkg/page"
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
)

type BlogCommentDao struct {
}

func (d BlogCommentDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table(blog.Comment{}.TableName())
}

// InsertComment 添加评论
func (d BlogCommentDao) InsertComment(comment blog.Comment) int64 {
	session := SqlDB.NewSession()
	session.Begin()

	insert, err := session.Insert(&comment)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()

	return insert
}

// DeleteComment 删除评论
func (d BlogCommentDao) DeleteComment(commentId uint) bool {
	session := SqlDB.NewSession()
	session.Begin()
	_, err := session.Where("id=?", commentId).Delete(&blog.Comment{})
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return false
	}
	session.Commit()
	return true
}

func (d BlogCommentDao) GetCommentById(commentId uint) *blog.Comment {
	var comment blog.Comment
	session := d.sql(SqlDB.NewSession())
	_, err := session.Where("id = ?", commentId).Get(&comment)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &comment
}

func (d BlogCommentDao) List(query blog.CommentQuery) (*[]blog.Comment, int64) {
	comment := make([]blog.Comment, 0)
	session := d.sql(SqlDB.NewSession())
	session.And("article_id=?", query.ArticleId)
	total, _ := page.GetTotal(session.Clone())
	err := session.Limit(query.PageSize, page.StartSize(query.PageNum, query.PageSize)).OrderBy("create_time DESC").Find(&comment)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil, 0
	}

	return &comment, total
}

func (d BlogCommentDao) GetCommentCountByArticle(articleId interface{}) int64 {

	session := d.sql(SqlDB.NewSession())

	if articleIds, ok := articleId.([]int64); ok {
		session.In("article_id", articleIds)
	} else {
		session.Where("article_id=?", articleId)
	}

	total, err := session.GroupBy("article_id").Count()

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return 0
	}

	return total
}
