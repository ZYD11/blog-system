package blog

import (
	"time"
)

type Article struct {
	Id         uint      `json:"id" xorm:"pk autoincr"`
	UserId     uint      `json:"user_id" xorm:"notnull"`
	CategoryId uint      `json:"category_id" xorm:"notnull"`
	Title      string    `json:"title" xorm:"varchar(50) notnull"`
	Content    string    `json:"content" xorm:"longtext notnull"`
	HeadImage  string    `json:"head_image" xorm:"longtext"`
	Status     uint      `json:"status" xorm:"tinyint(4) default(0)"`
	CreatedAt  time.Time `json:"created_at" xorm:"datetime"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"datetime"`
}

type ArticleInfo struct {
	Id         string `json:"id"`
	CategoryId uint   `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	HeadImage  string `json:"head_image"`
	CreatedAt  Time   `json:"created_at"`
}

type ArticleQuery struct {
	BeginTime  string `form:"beginTime"` //开始时间
	EndTime    string `form:"endTime"`   //结束时间
	PageNum    int    `form:"pageNum"`   //当前页码
	PageSize   int    `form:"pageSize"`  //显示条数
	CategoryId string `form:"categoryId"`
	Keyword    string `form:"keyword"`
}

func (receiver Article) TableName() string {
	return "blog_article"
}
