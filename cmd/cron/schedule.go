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
        id      uint64
        jobName string

        method  interface{}
        params  string  //多个参数 都改为
}

func (s *schedule) Run() {
        if s.method == nil {
                panic("Invalid Parameters")
        }
        logger.Infof("cron exec: %d %s runing...", s.id, s.jobName)
        if v, ok := s.method.(func()); ok {
                v()
        } else {
                v := s.method.(func(...interface{}))
                strs := strings.Split(s.params, ",")
                sList := make([]interface{}, len(strs))
                for k := range strs {
                        sList[k] = strs[k]
                }
                v(sList...)
        }
}
