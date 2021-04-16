package datamodels

//@Title sys_oper_log_.go
//@Description 后台系统操作日志
//@Author bertang
//@Created bertang 2021/1/11 11:47 上午

import "time"

//SysOperLog 系统操作日志
type SysOperLog struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;comment:id主键"`
	CreatedAt     time.Time //创建时间，同时是操作时间
	Path          string    `gorm:"type:varchar(32);not null;default:'';comment:请求路径"`
	Type          int8      `gorm:"type:tinyint(1);not null;default:0;comment:操作类型，0-其他，1-新增,2-修改,3-删除,4-上传"`
	RequestMethod string    `gorm:"type:char(4);not null;default:'';comment:请求方式"`
	OperName      string    `gorm:"type:varchar(32);not null;default:'';comment:操作员"`
	OperID        uint      `gorm:"not nulll;default:0;comment:操作员id"`
	Params        string    `gorm:"type:varchar(512);not null;default:'';comment:请求参数"`
	Result        string    `gorm:"type:varchar(512);not null;default:'';comment:请求结果json"`
}
