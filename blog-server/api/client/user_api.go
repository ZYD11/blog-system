package client

import (
	"blog-server/models/blog"
	"blog-server/service"
	"blog-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserApi 用户操作api
type UserApi struct {
	userService service.BlogUserService
	roleService service.RoleService
	postService service.PostService
}

// 注册
func (userApi UserApi) Register(ctx *gin.Context) {

	// 获取参数
	registerBody := blog.User{}
	if ctx.BindJSON(&registerBody) == nil {
		code, msg := userApi.userService.Register(registerBody.UserName, registerBody.PhoneNumber, registerBody.Password)
		if code == http.StatusOK {
			util.Success(ctx, nil, msg)
		} else {
			util.Fail(ctx, nil, msg)
		}
	} else {
		util.Fail(ctx, nil, "参数绑定错误")
	}
}

// 登录
func (userApi UserApi) Login(ctx *gin.Context) {
	// 获取参数
	loginBody := blog.User{}
	if ctx.BindJSON(&loginBody) == nil {

		code, msg, token := userApi.userService.Login(loginBody.PhoneNumber, loginBody.Password)
		if code == http.StatusOK {
			//将token存入到redis中
			// user_util.SaveRedisToken(loginBody.PhoneNumber, token)
			util.Success(ctx, gin.H{"token": token}, msg)
		} else {
			util.Fail(ctx, nil, msg)
		}
	} else {
		util.Fail(ctx, nil, "参数绑定错误")
	}
}

func (userApi UserApi) LogOut(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	if user != nil {
		ctx.Set("user", "")
	}

	util.Success(ctx, nil, "已注销！")
}

// 获取登录用户信息
func (userApi UserApi) GetUserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	util.Success(ctx, gin.H{"id": user.(blog.User).Id, "avatar": user.(blog.User).Avatar}, "登录获取信息成功")
}

// 获取简要信息
func (userApi UserApi) GetBriefInfo(ctx *gin.Context) {
	userId := ctx.Params.ByName("id")
	user, _ := ctx.Get("user")

	var curUser blog.User
	if userId == strconv.Itoa(int(user.(blog.User).Id)) {
		curUser = user.(blog.User)
	} else {
		paramUserId, _ := strconv.Atoi(userId)
		curUser = *(userApi.userService.GetUserByUserId(uint(paramUserId)))

		if curUser.Id == 0 {
			util.Fail(ctx, nil, "用户不存在")
			return
		}
	}

	util.Success(ctx, gin.H{"id": curUser.Id, "name": curUser.UserName, "avatar": curUser.Avatar, "loginId": user.(blog.User).Id}, "查找成功")

}

// 获取详细信息
func (userApi UserApi) GetDetailedInfo(ctx *gin.Context) {
	// 获取path中的userId
	userId := ctx.Params.ByName("id")
	paramUserId, _ := strconv.Atoi(userId)
	// 判断用户身份
	user, _ := ctx.Get("user")
	//var self bool
	var curUser blog.User
	if userId == strconv.Itoa(int(user.(blog.User).Id)) {
		//self = true
		curUser = user.(blog.User)
	} else {
		//self = false
		curUser = *(userApi.userService.GetUserByUserId(uint(paramUserId)))

		if curUser.Id == 0 {
			util.Fail(ctx, nil, "用户不存在")
			return
		}
	}
	// 返回用户详细信息
	var articles, collects *[]blog.ArticleInfo
	var following *[]blog.UserInfo
	var collist, follist []string
	collist = ToStringArray(curUser.Collects)
	follist = ToStringArray(curUser.Following)

	articleApi := ArticleApi{}
	articles = articleApi.articleService.GetArticlesByUserId(paramUserId)
	collects = articleApi.articleService.GetArticlesByIds(collist)
	following = userApi.userService.GetFollowingByUserIds(follist)
	util.Success(ctx, gin.H{"id": curUser.Id, "name": curUser.UserName, "avatar": curUser.Avatar, "loginId": user.(blog.User).Id, "articles": articles, "collects": collects, "following": following, "fans": curUser.Fans}, "查找成功")
}

// 修改头像
func (userApi UserApi) ModifyAvatar(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取参数
	var requestUser blog.User
	ctx.Bind(&requestUser)
	avatar := requestUser.Avatar
	// 查找用户
	var curUser blog.User
	curUser = *(userApi.userService.GetUserByUserId(user.(blog.User).Id))
	curUser.Avatar = avatar

	// 更新信息
	if userApi.userService.Update(curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}

	util.Success(ctx, nil, "更新成功")
}

// 修改名称
func (userApi UserApi) ModifyName(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取参数
	var requestUser blog.User
	ctx.Bind(&requestUser)
	username := requestUser.UserName
	// 查找用户
	var curUser blog.User
	curUser = *(userApi.userService.GetUserByUserId(user.(blog.User).Id))
	curUser.UserName = username

	// 更新信息
	if userApi.userService.Update(curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}

	util.Success(ctx, nil, "更新成功")
}

// 查询收藏
func (userApi UserApi) Collects(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	id := ctx.Params.ByName("id")
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)

	// 判断是否已收藏
	for i := 0; i < len(curUser.Collects); i++ {
		if curUser.Collects[i] == id {
			util.Success(ctx, gin.H{"collected": true, "index": i}, "查询成功")
			return
		}
	}
	util.Success(ctx, gin.H{"collected": false}, "查询成功")
}

// 增加收藏
func (userApi UserApi) NewCollect(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// 查找用户
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)
	var newCollects []string
	newCollects = append(curUser.Collects, id)
	curUser.Collects = newCollects

	// 更新收藏夹
	if userApi.userService.Update(*curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}
	util.Success(ctx, nil, "更新成功")
}

// 取消收藏
func (userApi UserApi) UnCollect(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的index
	index, _ := strconv.Atoi(ctx.Params.ByName("index"))
	// 查找用户
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)
	var newCollects []string
	newCollects = append(curUser.Collects[:index], curUser.Collects[index+1:]...)
	curUser.Collects = newCollects

	// 更新收藏夹
	if userApi.userService.Update(*curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}
	util.Success(ctx, nil, "更新成功")
}

// 查询关注
func (userApi UserApi) Following(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	id := ctx.Params.ByName("id")
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)

	// 判断是否已关注
	for i := 0; i < len(curUser.Following); i++ {
		if curUser.Following[i] == id {
			util.Success(ctx, gin.H{"followed": true, "index": i}, "查询成功")
			return
		}
	}
	util.Success(ctx, gin.H{"collected": false}, "查询成功")
}

// 添加关注
func (userApi UserApi) NewFollow(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// 查找用户
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)

	newFollowing := append(curUser.Following, id)
	curUser.Following = newFollowing

	// 更新粉丝
	if userApi.userService.Update(*curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}

	// 更新粉丝数
	var followUser *blog.User
	paramId, _ := strconv.Atoi(id)
	followUser = userApi.userService.GetUserByUserId(uint(paramId))
	followUser.Fans = followUser.Fans + 1

	if userApi.userService.Update(*followUser) <= 0 {
		util.Fail(ctx, nil, "更新粉丝数量失败！")
		return
	}
	util.Success(ctx, nil, "更新成功")
}

// 取消关注
func (userApi UserApi) UnFollow(ctx *gin.Context) {
	// 获取用户ID
	user, _ := ctx.Get("user")
	// 获取path中的index
	index, _ := strconv.Atoi(ctx.Params.ByName("index"))

	// 查找用户
	var curUser *blog.User
	curUser = userApi.userService.GetUserByUserId(user.(blog.User).Id)

	newFollowing := append(curUser.Following[:index], curUser.Following[index+1:]...)
	followId := curUser.Following[index]
	curUser.Following = newFollowing

	// 更新粉丝
	if userApi.userService.Update(*curUser) <= 0 {
		util.Fail(ctx, nil, "更新失败")
		return
	}

	// 更新粉丝数
	var followUser *blog.User
	paramId, _ := strconv.Atoi(followId)
	followUser = userApi.userService.GetUserByUserId(uint(paramId))
	followUser.Fans = followUser.Fans - 1

	if userApi.userService.Update(*followUser) <= 0 {
		util.Fail(ctx, nil, "更新粉丝数量失败！")
		return
	}
	util.Success(ctx, nil, "更新成功")
}

// 将自定义类型转化为字符串数组
func ToStringArray(l []string) (a blog.Array) {
	for i := 0; i < len(a); i++ {
		l = append(l, a[i])
	}
	return l
}
