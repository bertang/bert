package cron

import (
	"errors"
	"bert/types"
	"github.com/robfig/cron"
	"sync"
)

var (
	ch chan interface{}
	jobs sync.Map
)

type ICronService interface {
	ActiveList()
}

type crontab struct {
	Func    interface{} `json:"-"`
	Comment string     `json:"comment"`
}

//注册定时任务
func Register(key string, funcName interface{})  {
	//只接收这两个类型的数据
	if _, ok:=funcName.(types.Func); !ok {
		if _, ok:= funcName.(types.FuncWithParams); !ok {
			panic("Invalid Parameter")
		}
	}
	if _, ok:= jobs.LoadOrStore(key, funcName); ok {
		panic(errors.New("Cron exits!"))
	}
}

func Start()  {
	c:=cron.New()
	c.Run()
}

func start() {
	ch = make(chan interface{})
	go func() {

	}()
}