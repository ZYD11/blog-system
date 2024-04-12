package dao

import (
	"blog-server/models/blog"
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
)

type BlogUserDao struct {
}

func (d BlogUserDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table(blog.User{}.TableName())
}

// 根据用条件查询用户数据
func (d BlogUserDao) GetUserByConditions(user blog.User) *blog.User {
	i, err := SqlDB.Get(&user)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	if i {
		return &user
	}
	return nil
}

// InsertUser 添加用户
func (d BlogUserDao) InsertUser(user blog.User) int64 {
	session := SqlDB.NewSession()
	session.Begin()

	insert, err := session.Insert(&user)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()

	return insert
}

func (d BlogUserDao) UpdateUser(user blog.User) int64 {
	session := SqlDB.NewSession()
	session.Begin()
	update, err := session.Where("id = ?", user.Id).AllCols().Update(&user)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()
	return update
}

func (d BlogUserDao) GetFollowingByUserIds(userIds []string) *[]blog.UserInfo {
	var users []blog.UserInfo
	session := d.sql(SqlDB.NewSession())
	err := session.Select("id, avatar, user_name").In("id", userIds).Find(&users)

	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
	}

	return &users
}
