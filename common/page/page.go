package page

import (
	"github.com/bertang/bert/common/config/application"
	"github.com/kataras/iris/v12"
)

//IPageHelper 分页帮助
type IPageHelper interface {
	Offset() int
	Limit() int
}

//Helper 分页助手
type Helper struct {
	Page   int   `json:"page"`
	Size   int   `json:"size"`
	Total  int64 `json:"total"`
	offset int
}

//Offset 获取offset
func (p *Helper) Offset() int {
	return (p.Page - 1) * p.Size
}

//Limit 获取limit
func (p *Helper) Limit() int {
	return p.Size
}

//NewPageHelper 生成新的分页助手
func NewPageHelper(ctx iris.Context) Helper {
	page := ctx.URLParamIntDefault("page", 1)
	size := ctx.URLParamIntDefault("size", 0)

	if size == 0 {
		size = application.GetAppConf().PageSize
	}

	helper := Helper{
		Page: page,
		Size: size,
	}

	return helper
}
