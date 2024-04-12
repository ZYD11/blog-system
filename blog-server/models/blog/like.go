package blog

import "time"

type Like struct {
	Id         uint      `json:"id" xorm:"pk autoincr"`
	ArticleId  uint      `json:"article_id" xorm:"notnull"`
	LikeUser   uint      `json:"like_user" xorm:"notnull"`
	CreateTime time.Time `json:"create_time" xorm:"datetime"`
}

func (receiver Like) TableName() string {
	return "blog_like"
}
