package view

//Response 响应结果
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//Result 输入所有的参数返回响应
//@code 响应码
//@message 响应消息
//@data 需要返回的数据
func Result(code int, message string, data interface{}) Response {
	return Response{code, message, data}
}
