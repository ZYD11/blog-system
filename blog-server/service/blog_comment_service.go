package service

import (
	"blog-server/dao"
	"blog-server/models/blog"
)

type BlogCommentService struct {
	commentDao dao.BlogCommentDao
}

func (s BlogCommentService) Add(comment blog.Comment) int64 {
	return s.commentDao.InsertComment(comment)
}

func (s BlogCommentService) Remove(commentId uint) bool {
	return s.commentDao.DeleteComment(commentId)
}

func (s BlogCommentService) GetCommentById(commentId uint) *blog.Comment {
	return s.commentDao.GetCommentById(commentId)
}

func (s BlogCommentService) Find(query blog.CommentQuery) (*[]blog.Comment, int64) {
	return s.commentDao.List(query)
}

func (s BlogCommentService) GetCommentCountByArticle(articleId interface{}) int64 {
	return s.commentDao.GetCommentCountByArticle(articleId)
}
