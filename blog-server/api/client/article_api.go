package client

import (
	"blog-server/models/blog"
	"blog-server/service"
	"blog-server/util"
	"blog-server/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// ArticleApi 博客文章操作api
type ArticleApi struct {
	articleService service.BlogArticleService
	commentService service.BlogCommentService
	likeService    service.BlogLikeService
	userService    service.BlogUserService
}

// 创建文章
func (articleApi ArticleApi) Create(ctx *gin.Context) {
	var articleRequest validate.CreateArticleRequest
	// 数据验证
	if err := ctx.ShouldBindJSON(&articleRequest); err != nil {
		util.Fail(ctx, nil, "数据错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	// 创建文章
	article := blog.Article{
		UserId:     user.(blog.User).Id,
		CategoryId: articleRequest.CategoryId,
		Title:      articleRequest.Title,
		Content:    articleRequest.Content,
		HeadImage:  articleRequest.HeadImage,
		CreatedAt:  time.Now(),
	}

	if articleApi.articleService.Add(article) <= 0 {
		util.Fail(ctx, nil, "发布失败")
		return
	}

	util.Success(ctx, gin.H{"id": article.Id}, "发布成功")
}

// 更新文章
func (articleApi ArticleApi) Update(ctx *gin.Context) {
	var articleRequest validate.CreateArticleRequest
	// 数据验证
	if err := ctx.ShouldBindJSON(&articleRequest); err != nil {
		util.Fail(ctx, nil, "数据错误")
		return
	}
	// 获取path中的id
	articleId := ctx.Params.ByName("id")
	// 查找文章
	var article *blog.Article
	article = articleApi.articleService.GetArticleById(articleId)

	if article == nil {
		util.Fail(ctx, nil, "文章不存在")
		return
	}

	user, _ := ctx.Get("user")
	userId := user.(blog.User).Id
	if userId != article.UserId {
		util.Fail(ctx, nil, "登录用户不正确")
		return
	}

	articleUpdate := blog.Article{
		Id:         article.Id,
		UserId:     userId,
		CategoryId: articleRequest.CategoryId,
		Title:      articleRequest.Title,
		Content:    articleRequest.Content,
		HeadImage:  articleRequest.HeadImage,
		UpdatedAt:  time.Now(),
	}

	fmt.Println(articleUpdate)
	// 更新文章
	if articleApi.articleService.Update(articleUpdate) <= 0 {
		util.Fail(ctx, nil, "修改失败")
		return
	}

	util.Success(ctx, nil, "修改成功")
}

// 删除文章
func (articleApi ArticleApi) Delete(ctx *gin.Context) {
	// 获取path中的id
	articleId := ctx.Params.ByName("id")

	// 查找文章
	var article *blog.Article
	article = articleApi.articleService.GetArticleById(articleId)
	if article == nil {
		util.Fail(ctx, nil, "文章不存在")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(blog.User).Id
	if userId != article.UserId {
		util.Fail(ctx, nil, "登录用户不正确")
		return
	}

	if articleApi.articleService.Remove(articleId) == false {
		util.Fail(ctx, nil, "删除失败")
		return
	}

	util.Success(ctx, nil, "删除成功")
}

// 文章展示
func (articleApi ArticleApi) Show(ctx *gin.Context) {
	// 获取path中的id
	articleId := ctx.Params.ByName("id")
	var article *blog.Article
	article = articleApi.articleService.GetArticleById(articleId)
	// 查找文章
	if article == nil {
		util.Fail(ctx, nil, "文章不存在")
		return
	}
	// 获取关键词、分类、分页参数
	util.Success(ctx, gin.H{"article": article}, "查找成功！")
}

// 文章展示列表
func (articleApi ArticleApi) List(ctx *gin.Context) {
	// 获取关键词、分类、分页参数
	keyword := ctx.DefaultQuery("keyword", "")
	categoryId := ctx.DefaultQuery("categoryId", "0")
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "5"))

	articleQuery := blog.ArticleQuery{
		PageNum:    pageNum,
		PageSize:   pageSize,
		Keyword:    keyword,
		CategoryId: categoryId,
	}
	// 页面内容
	var articles *[]blog.Article
	// 文章总数
	var count int64

	articleCount := make(map[string]map[string]int64)
	articles, count = articleApi.articleService.Find(articleQuery)

	var userIds []string
	articleUsers := make(map[uint]blog.UserInfo)

	for _, article := range *articles {
		countMap := make(map[string]int64)
		comment_count := articleApi.commentService.GetCommentCountByArticle(article.Id)
		like_count := articleApi.likeService.GetLikeCountByArticle(article.Id)
		key := strconv.Itoa(int(article.Id))
		countMap["like_count"] = like_count
		countMap["comment_count"] = comment_count
		articleCount[key] = countMap
		userIds = append(userIds, strconv.Itoa(int(article.UserId)))
	}

	if len(userIds) > 0 {
		userData := articleApi.userService.GetFollowingByUserIds(userIds)
		if userData != nil {
			for _, userList := range *userData {
				fmt.Println(userList)
				articleUsers[userList.Id] = userList
			}
		}
	}
	util.Success(ctx, gin.H{"article": articles, "articleCount": articleCount, "articleUsers": articleUsers, "count": count}, "查找成功")
}

// 评论列表
func (articleApi ArticleApi) CommentList(ctx *gin.Context) {

	articleId := ctx.DefaultQuery("articleId", "0")
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "5"))

	commentQuery := blog.CommentQuery{
		PageNum:   pageNum,
		PageSize:  pageSize,
		ArticleId: articleId,
	}
	// 页面内容
	var comments *[]blog.Comment

	// 评论总数
	var count int64
	comments, count = articleApi.commentService.Find(commentQuery)

	var userIds []string
	commentUsers := make(map[uint]blog.UserInfo)
	for _, comment := range *comments {
		userIds = append(userIds, strconv.Itoa(int(comment.Reviewer)))
	}

	if len(userIds) > 0 {
		var userData = articleApi.userService.GetFollowingByUserIds(userIds)
		if userData != nil {
			for _, userList := range *userData {
				commentUsers[userList.Id] = userList
			}
		}
	}

	util.Success(ctx, gin.H{"comments": comments, "commentUsers": commentUsers, "count": count}, "查找成功")
}

// 新增评论
func (articleApi ArticleApi) NewComment(ctx *gin.Context) {
	var commentRequest validate.CreateCommentRequest
	// 数据验证
	if err := ctx.ShouldBindJSON(&commentRequest); err != nil {
		fmt.Println(err)
		util.Fail(ctx, nil, "数据错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	comment := blog.Comment{
		ArticleId:  commentRequest.ArticleId,
		Content:    commentRequest.Content,
		Reviewer:   user.(blog.User).Id,
		CreateTime: time.Now(),
	}

	if articleApi.commentService.Add(comment) <= 0 {
		util.Fail(ctx, nil, "发表失败")
		return
	}

	util.Success(ctx, gin.H{"id": comment.Id}, "发表成功")
}

// 删除评论
func (articleApi ArticleApi) DeleteComment(ctx *gin.Context) {
	fmt.Println(ctx)
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的index
	commentId, _ := strconv.Atoi(ctx.Params.ByName("id"))

	// 查找文章
	var comment *blog.Comment
	comment = articleApi.commentService.GetCommentById(uint(commentId))
	if comment == nil {
		util.Fail(ctx, nil, "评论不存在")
		return
	}

	// 获取登录用户
	userId := user.(blog.User).Id

	if userId != comment.Reviewer {
		util.Fail(ctx, nil, "无权限删除他人评论")
		return
	}

	if articleApi.commentService.Remove(uint(commentId)) == false {
		util.Fail(ctx, nil, "删除失败")
		return
	}

	util.Success(ctx, nil, "删除成功")
}

// 新增点赞
func (articleApi ArticleApi) NewLike(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Params.ByName("id"))
	// 数据验证
	if err != nil || articleId <= 0 {
		util.Fail(ctx, nil, "数据错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	like := blog.Like{
		ArticleId:  uint(articleId),
		LikeUser:   user.(blog.User).Id,
		CreateTime: time.Now(),
	}

	if articleApi.likeService.Add(like) <= 0 {
		util.Fail(ctx, nil, "发表失败")
		return
	}

	util.Success(ctx, gin.H{"id": like.Id}, "发表成功")
}

// 取消点赞
func (articleApi ArticleApi) UnLike(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	likeId := ctx.Params.ByName("id")

	// 查找点赞id
	var like *blog.Like
	like = articleApi.likeService.GetLikeById(likeId)
	if like == nil {
		util.Fail(ctx, nil, "点赞不存在！")
		return
	}

	// 获取登录用户
	userId := user.(blog.User).Id
	if userId != like.LikeUser {
		util.Fail(ctx, nil, "无权限取消他人点赞")
		return
	}

	if articleApi.likeService.Remove(likeId) == false {
		util.Fail(ctx, nil, "取消失败")
		return
	}

	util.Success(ctx, nil, "取消成功")
}

func (articleApi ArticleApi) Like(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	id, _ := strconv.Atoi(ctx.Params.ByName("id"))

	currentLike := blog.Like{
		LikeUser:  user.(blog.User).Id,
		ArticleId: uint(id),
	}

	like := articleApi.likeService.GetLikeByCondition(currentLike)

	if like != nil {
		util.Success(ctx, gin.H{"liked": true, "like_id": like.Id}, "查询成功")
		return
	}

	util.Success(ctx, gin.H{"liked": false}, "查询成功")
}
