package service

import (
	"blog-server/dao"
	"blog-server/models/blog"
	"blog-server/pkg/common"
	"fmt"
	"github.com/druidcaesa/gotool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type BlogUserService struct {
	userDao dao.BlogUserDao
}

// 用户注册业务处理
func (s BlogUserService) Register(userName string, phoneNumber string, password string) (code int, msg string) {
	user := s.GetUserByPhoneNumber(phoneNumber)
	fmt.Println(user)
	if user != nil && user.Id != 0 {
		return 422, "用户已存在"
	}

	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 注册成功，保存数据到数据库
	newUser := blog.User{
		UserName:    userName,
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
		Avatar:      "/images/default_avatar.png",
		Collects:    blog.Array{},
		Following:   blog.Array{},
		Fans:        0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	//添加用户数据库操作
	insertUser := s.userDao.InsertUser(newUser)

	if insertUser > 0 {
		return 200, "注册成功！"
	} else {
		return 500, "注册失败！"
	}

}

// Login 用户登录业务处理
func (s BlogUserService) Login(phoneNumber string, password string) (code int, msg string, token string) {
	user := s.GetUserByPhoneNumber(phoneNumber)

	if user == nil {
		return 422, "用户不存在", ""
	}
	// 对比密码，判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 422, "密码错误", ""
	}

	//生成token
	token, err := common.GenerateToken(user.Id)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return 500, "系统异常", ""
	}

	return 200, "登录成功！", token
}

func (s BlogUserService) Add(newUser blog.User) int64 {
	//添加用户数据库操作
	return s.userDao.InsertUser(newUser)
}

func (s BlogUserService) Update(newUser blog.User) int64 {
	//更新用户数据库操作
	return s.userDao.UpdateUser(newUser)
}
func (s BlogUserService) GetUserByUserName(username string) *blog.User {
	user := blog.User{}
	user.UserName = username
	return s.userDao.GetUserByConditions(user)
}

func (s BlogUserService) GetUserByPhoneNumber(phoneNumber string) *blog.User {
	user := blog.User{}
	user.PhoneNumber = phoneNumber
	return s.userDao.GetUserByConditions(user)
}

func (s BlogUserService) GetUserByUserId(userId uint) *blog.User {
	user := blog.User{}
	user.Id = userId
	return s.userDao.GetUserByConditions(user)
}

func (s BlogUserService) GetFollowingByUserIds(userId []string) *[]blog.UserInfo {
	return s.userDao.GetFollowingByUserIds(userId)
}
