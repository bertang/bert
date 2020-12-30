package page

import (
	"github.com/bertang/bert/common/config/application"
	"github.com/kataras/iris/v12"
)

//Helper 分页助手
type Helper struct {
	Page   int   `json:"page"`
	Size   int   `json:"size"`
	Total  int64 `json:"total"`
	offset int
}

//Offset 获取offset
func (p *Helper) Offset() int {
	return p.offset
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
	helper.offset = (helper.Page - 1) * helper.Size
	return helper
}
