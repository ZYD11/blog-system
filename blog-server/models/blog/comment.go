package blog

import "time"

type Comment struct {
	Id         uint      `json:"id" xorm:"pk autoincr"`
	ArticleId  uint      `json:"article_id" xorm:"notnull"`
	Content    string    `json:"content" xorm:"longtext notnull"`
	Reviewer   uint      `json:"reviewer" xorm:"notnull"`
	Status     uint      `json:"status" xorm:"tinyint default(0)"`
	CreateTime time.Time `json:"create_time" xorm:"datetime"`
	UpdateTime time.Time `json:"update_time" xorm:"datetime"`
}

type CommentQuery struct {
	PageNum   int    `form:"pageNum"`  //当前页码
	PageSize  int    `form:"pageSize"` //显示条数
	ArticleId string `form:"articleId"`
}

func (receiver Comment) TableName() string {
	return "blog_comment"
}
