package services

import "bert/datamodels"

type ICrontabService interface {
	//获取已正常激活的定时任务
	JobActived() ([]*datamodels.SysJob, error)
	//分页获取所有的定时任务
	List(int,int)([]*datamodels.SysJob,int64, error)

}
