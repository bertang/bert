package controller

import (
	"bytes"
	"net/http"
	"time"

	"github.com/bertang/bert/common/consts"
	"github.com/bertang/bert/common/jwt"
	jsoniter "github.com/json-iterator/go"

	"github.com/bertang/bert/apis/view"
	"github.com/bertang/bert/common/hmessage"
	"github.com/bertang/bert/datamodels"
	"github.com/bertang/bert/services"
	"github.com/kataras/iris/v12"
	jwt2 "github.com/kataras/iris/v12/middleware/jwt"
)

//BaseController 基础控制器
type BaseController struct {
	Ctx            iris.Context
	OperLogService services.ISysOperLogService
}

//Success 返回成功
func (b *BaseController) Success() view.Response {
	return b.Result(200, hmessage.SUCCESS, nil)
}

//SuccessWithData 有Data返回
//@data 需要返回的数据
func (b *BaseController) SuccessWithData(data interface{}) view.Response {
	return b.Result(http.StatusOK, hmessage.SUCCESS, data)
}

//Result 最终返回
func (b *BaseController) Result(code int, message string, data interface{}) view.Response {
	ret := view.Response{Code: code, Message: message, Data: data}
	if b.Ctx.Method() == iris.MethodGet {
		return ret
	}
	if b.OperLogService != nil {
		logModel := &datamodels.SysOperLog{
			CreatedAt:     time.Now(),
			Path:          b.Ctx.Path(),
			RequestMethod: b.Ctx.Method(),
			OperName:      "",
			OperID:        0,
			Params:        "",
			Result:        "",
		}

		//如果是其他格式，如果是上传不记录
		if b.Ctx.GetHeader("Content-Type") != "application/json" {
			logModel.Type = consts.OperLogTypeOther
		} else {
			if b.Ctx.Method() == iris.MethodPost {
				logModel.Type = consts.OperLogTypeAdd
			} else if b.Ctx.Method() == iris.MethodDelete {
				logModel.Type = consts.OperLogTypeDelete
			} else if b.Ctx.Method() == iris.MethodPut {
				logModel.Type = consts.OperLogTypeUpdate
			}
			p, _ := b.Ctx.GetBody()
			logModel.Params = string(bytes.ReplaceAll(bytes.ReplaceAll(p, []byte{'\n'}, []byte{}), []byte("    "), []byte{}))
		}

		userClaim := jwt2.Get(b.Ctx).(*jwt.UserClaim)
		logModel.OperName = userClaim.Name
		logModel.OperID = userClaim.ID
		retStr, _ := jsoniter.Marshal(ret)
		logModel.Result = string(retStr)
		_ = b.OperLogService.Add(logModel)
	}
	return ret
}

//ServerError 返回服务器内部错误
func (b *BaseController) ServerError() view.Response {
	return b.Result(http.StatusInternalServerError, hmessage.CodeMessage[http.StatusInternalServerError], nil)
}

//ParameterError 返回参数错误，主要用于没有读取到json值时
func (b *BaseController) ParameterError() view.Response {
	return b.Result(hmessage.ErrParamsCode, hmessage.CodeMessage[hmessage.ErrParamsCode], nil)
}

//ValidateError 返回难错误
func (b *BaseController) ValidateError(data interface{}) view.Response {
	return b.Result(hmessage.ErrParamsCode, hmessage.CodeMessage[hmessage.ErrParamsCode], data)
}

//NotFound 资源不存在
func (b *BaseController) NotFound() view.Response {
	return b.Result(http.StatusNotFound, hmessage.CodeMessage[http.StatusNotFound], nil)
}

//Forbidden 无权限访问
func (b *BaseController) Forbidden(err error) view.Response {
	return b.Result(http.StatusForbidden, err.Error(), nil)
}

//UploadError 文件上传失败错误
func (b *BaseController) UploadError() view.Response {
	return b.Result(hmessage.ErrParamsCode, "文件上传失败", nil)
}
