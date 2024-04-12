package service

import (
	"blog-server/dao"
	"blog-server/models/blog"
)

type BlogLikeService struct {
	likeDao dao.BlogLikeDao
}

// 新增点赞
func (s BlogLikeService) Add(like blog.Like) int64 {
	return s.likeDao.InsertLike(like)
}

// 取消点赞
func (s BlogLikeService) Remove(likeId string) bool {
	return s.likeDao.DeleteLike(likeId)
}

func (s BlogLikeService) GetLikeById(likeId string) *blog.Like {
	return s.likeDao.GetLikeById(likeId)
}

func (s BlogLikeService) GetLikeByCondition(like blog.Like) *blog.Like {
	return s.likeDao.GetLikeByCondition(like)
}

func (s BlogLikeService) GetLikeCountByArticle(articleId interface{}) int64 {
	return s.likeDao.GetLikeCountByArticle(articleId)
}
