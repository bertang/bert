package hmessage

import "net/http"

const (
	//SUCCESS 操作成功
	SUCCESS = "操作成功"
	//UNAUTHORIZED 未授权
	UNAUTHORIZED = "未授权的资源"
	//FORBIDDEN 没有权限
	FORBIDDEN = "没有权限，请联系管理员授权"
	//NOT_FOUND 资源不存在
	NOT_FOUND = "资源不存在"
	//SERVER_ERROR 系统内部错误
	SERVER_ERROR = "系统内部错误"
)

var (
	CodeMessage = map[int]string{
		http.StatusOK:                  SUCCESS,
		http.StatusUnauthorized:        UNAUTHORIZED,
		http.StatusForbidden:           FORBIDDEN,
		http.StatusNotFound:            NOT_FOUND,
		http.StatusInternalServerError: SERVER_ERROR,
	}
)
