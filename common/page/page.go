package page

import "github.com/kataras/iris/v12"

type IPageHelper interface {
	//Offset 获取offset
	Offset() int
	//Limit 获取limit
	Limit() int
}
type pageHelper struct {
	Page   int `json:"page"`
	Size   int `json:"size"`
	Total  int `json:"total"`
	offset int
}

//获取offset
func (p *pageHelper) Offset() int {
	return p.offset
}

//获取limit
func (p *pageHelper) Limit() int {
	return p.Size
}

func NewPageHelper(ctx iris.Context) IPageHelper {
	page := ctx.URLParamIntDefault("page", 1)
	size := ctx.URLParamIntDefault("size", 20)

	helper := &pageHelper{
		Page: page,
		Size: size,
	}
	helper.Size = (helper.Page - 1) * helper.Size
	return helper
}
