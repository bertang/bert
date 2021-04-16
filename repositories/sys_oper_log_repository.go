package repositories

import "gorm.io/gorm"

//@file sys_oper_log_repository
//@Description 系统操作日志
//@Author bertang
//@Created 2021/01/12 11:00:00

//ISysOperLogRepository 系统操作日志
type ISysOperLogRepository interface {
	IBaseRepo
}
type sysOperLogRepository struct {
	BaseRepository
}

//NewSysOperLog 新建
func NewSysOperLog(db *gorm.DB) ISysOperLogRepository {
	return &sysOperLogRepository{BaseRepository: BaseRepository{Db: db}}
}
