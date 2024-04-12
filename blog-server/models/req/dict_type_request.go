package req

import "blog-server/pkg/base"

type DictTypeQuery struct {
	base.GlobalQuery
	DictName string `form:"dictName"`
	Status   string `form:"status"`
	DictType string `form:"dictType"`
}
