package blog

import "time"

type User struct {
	Id          uint      `xorm:"pk autoincr" json:"user_id"`
	UserName    string    `xorm:"varchar(20) notnull"`
	PhoneNumber string    `xorm:"varchar(20) notnull unique"`
	Password    string    `xorm:"varchar(255) notnull"`
	Avatar      string    `xorm:"varchar(255) notnull"`
	List        Array     `xorm:"longtext"`
	Collects    Array     `xorm:"longtext"`
	Following   Array     `xorm:"longtext"`
	Fans        int       `xorm:"int(11)"`
	CreatedAt   time.Time `xorm:"datetime"`
	UpdatedAt   time.Time `xorm:"datetime"`
}

type UserInfo struct {
	Id       uint   `json:"user_id"`
	Avatar   string `json:"avatar"`
	UserName string `json:"userName"`
}

func (receiver User) TableName() string {
	return "blog_users"
}
