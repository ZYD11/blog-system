package blog

type Category struct {
	Id           uint   `json:"id" xorm:"char(36) pk"`
	CategoryName string `json:"name" xorm:"varchar(50) notnull"`
}

func (receiver Category) TableName() string {
	return "blog_category"
}
