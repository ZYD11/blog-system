package admin

import (
	"blog-server/models/blog"
	"blog-server/pkg/resp"
	"blog-server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ArticleApi struct {
	articleService service.BlogArticleService
	userService    service.BlogUserService
}

// List 查询设置列表
func (a ArticleApi) List(c *gin.Context) {
	query := blog.ArticleQuery{}
	if c.Bind(&query) != nil {
		resp.ParamError(c)
		return
	}
	find, i := a.articleService.Find(query)
	var userIds []string
	publishUsers := make(map[uint]blog.UserInfo)
	for _, findData := range *find {
		userIds = append(userIds, strconv.Itoa(int(findData.UserId)))
	}

	if len(userIds) > 0 {
		var userData = a.userService.GetFollowingByUserIds(userIds)
		if userData != nil {
			for _, userList := range *userData {
				publishUsers[userList.Id] = userList
			}
		}
	}

	resp.OK(c, gin.H{
		"list":          find,
		"total":         i,
		"size":          query.PageSize,
		"publish_users": publishUsers,
	})
}

// Get 查询数据
func (a ArticleApi) Get(c *gin.Context) {
	articleId := c.Param("id")
	resp.OK(c, a.articleService.GetArticleById(articleId))
}

// 文章审核
func (a ArticleApi) Audit(c *gin.Context) {
	var postData map[string]interface{}
	if c.Bind(&postData) != nil {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
		return
	}
	articleId, ok := postData["articleId"]

	if !ok {
		c.JSON(500, resp.ErrorResp("参数错误"))
		return
	}

	fmt.Println(articleId)

	//if a.articleService.BathAudit(strconv.Itoa(articleId)) > 0 {
	//	c.JSON(200, resp.Success("审核成功"))
	//} else {
	//	c.JSON(500, resp.Success("审核失败"))
	//}
}

// Delete 删除数据
func (a ArticleApi) Delete(c *gin.Context) {
	ids := c.Param("ids")

	if a.articleService.Remove(ids) {
		c.JSON(200, resp.Success("删除成功"))
	} else {
		c.JSON(500, resp.Success("删除失败"))
	}
}
