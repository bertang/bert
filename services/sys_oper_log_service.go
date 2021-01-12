package services

//@Title sys_oper_log_service.go
//@Description
//@Author bertang
//@Created bertang 2021/1/11 2:31 下午

import (
	"github.com/bertang/bert/datamodels"
	"github.com/bertang/bert/repositories"
)

//ISysOperLogService 日志服务
type ISysOperLogService interface {
	Add(model *datamodels.SysOperLog) error
}

type sysOperLogService struct {
	repo repositories.ISysOperLogRepository
}

//Add 添加日志
func (s *sysOperLogService) Add(model *datamodels.SysOperLog) error {
	return s.repo.Add(model)
}

//NewSysOperLogService 创建日志服务
func NewSysOperLogService(repo repositories.ISysOperLogRepository) ISysOperLogService {
	return &sysOperLogService{repo: repo}
}
