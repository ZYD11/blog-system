package dao

import (
	"blog-server/models/blog"
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
)

type BlogLikeDao struct {
}

func (d BlogLikeDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table(blog.Like{}.TableName())
}

// 新增点赞
func (d BlogLikeDao) InsertLike(like blog.Like) int64 {
	session := SqlDB.NewSession()
	session.Begin()

	insert, err := session.Insert(&like)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()

	return insert
}

// 取消点赞
func (d BlogLikeDao) DeleteLike(likeId string) bool {
	session := SqlDB.NewSession()
	session.Begin()
	_, err := session.Where("id=?", likeId).Delete(&blog.Like{})
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return false
	}
	session.Commit()
	return true
}

func (d BlogLikeDao) GetLikeById(likeId string) *blog.Like {
	var like blog.Like
	session := d.sql(SqlDB.NewSession())
	_, err := session.Where("id = ?", likeId).Get(&like)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &like
}

func (d BlogLikeDao) GetLikeByCondition(like blog.Like) *blog.Like {
	i, err := SqlDB.Get(&like)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	if i {
		return &like
	}
	return nil
}

func (d BlogLikeDao) GetLikeCountByArticle(articleId interface{}) int64 {

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
