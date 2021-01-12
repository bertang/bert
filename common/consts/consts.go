//@Title consts.go
//@Description 
//@Author bertang
//@Created bertang 2021/1/12 11:18 上午
package consts


//操作日志相关
const (
    OperLogTypeOther = iota //其他
    OperLogTypeAdd //新增
    OperLogTypeUpdate //更新
    OperLogTypeDelete //删除
    OperLogTypeUpload //上传
)