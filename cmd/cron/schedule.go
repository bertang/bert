package cron

import (
	"strings"

	"github.com/bertang/bert/common/logger"
)

//Callback 有参函数
type jobFuncParams func(...interface{})

//JobFunc 无参函数
type jobFunc func()

type schedule struct {
	id      uint //数据库中保存的id
	jobName string //数据库中定义的定时任务名 字
	method interface{} //需要执行的函数
	params string //多个参数 都改为
}

func (s *schedule) Run() {
	if s.method == nil {
		panic("Invalid Parameters")
	}

	logger.Infof("cron exec: %d %s runing...", s.id, s.jobName)
	//无参函数执行
	if v, ok := s.method.(jobFunc); ok {
		v()
		return
	}

	//有参函数执行
	v := s.method.(jobFuncParams)
	//分解参数，拼接
	pStr := strings.Split(s.params, ",")
	sList := make([]interface{}, len(pStr))
	for k := range pStr {
		sList[k] = pStr[k]
	}

	//执行
	v(sList...)
}
