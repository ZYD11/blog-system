package service

import (
	"github.com/druidcaesa/gotool"
	"github.com/gin-gonic/gin"
	"blog-server/models/response"
	"blog-server/pkg/jwt"
	"strings"
)

type LoginService struct {
	userService UserService
}

// Login 用户登录业务处理
func (s LoginService) Login(name string, password string) (bool, string) {
	user := s.userService.GetUserByUserName(name)
	if user == nil {
		return false, "用户不存在"
	}
	if !gotool.BcryptUtils.CompareHash(user.Password, password) {
		return false, "密码错误"
	}
	//生成token
	token, err := jwt.NewJWT().CreateUserToken(s.userService.GetUserById(user.UserId))
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return false, ""
	}
	//数据存储到redis中
	return true, token
}

// LoginUser 获取当前登录用户
func (s LoginService) LoginUser(c *gin.Context) *response.UserResponse {
	token := c.Request.Header.Get("Authorization")
	str := strings.Split(token, " ")
	j := jwt.NewJWT()
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(str[1])
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
	}
	info := claims.UserInfo
	return &info
}
