package cron

import (
	"errors"
	"reflect"
	"sync"

	"github.com/bertang/bert/common/logger"
	"github.com/robfig/cron"
)

var (
	ch      chan interface{}
	jobs    sync.Map
	service ICronService
	mIndex  = new(modelIndex)
)

type modelIndex struct {
	name       int
	express    int
	id         int
	parameters int
	key        int
}

//ICronService 获取定时任务结构体
type ICronService interface {
	JobList() interface{}
}

//SetService 设置获取定时任务的service 以便定制化的数据库设计
func SetService(s ICronService) {
	service = s
}

type crontab struct {
	Func    interface{} `json:"-"`
	Comment string      `json:"comment"`
}

//Register 注册定时任务 由于不支持动态的获取函数方法。所有采用注册制
func Register(key string, funcName interface{}, comment string) {
	//只接收这两个类型的数据
	if _, ok := funcName.(func()); !ok {
		if _, ok := funcName.(func(...interface{})); !ok {
			panic("Invalid Parameter")
		}
	}
	if _, ok := jobs.LoadOrStore(key, &crontab{Func: funcName, Comment: comment}); ok {
		panic(errors.New("Cron Exits"))
	}
}

//Start 开始执行定时任务
func Start() {
	start()
}

func start() {
	ch = make(chan interface{})
	go func() {
		if service == nil {
			return
		}

		jobList := service.JobList()
		if jobList == nil {
			return
		}

		//判断传入类型
		jobListType := reflect.TypeOf(jobList)
		jobListRv := reflect.ValueOf(jobList)
		if jobListType.Kind() == reflect.Ptr {
			jobListRv = jobListRv.Elem()
			jobListType = jobListType.Elem()
		}

		if jobListType.Kind() != reflect.Slice {
			return
		}

		child := jobListType.Elem()
		if child.Kind() == reflect.Ptr {
			child = child.Elem()
		}
		//获取tag的位置
		for k := 0; k < child.NumField(); k++ {
			name := child.Field(k).Tag.Get("cron_key")
			switch name {
			case "name":
				mIndex.name = k
			case "express":
				mIndex.express = k

			case "params":
				mIndex.parameters = k
			case "key":
				mIndex.key = k
			}
			if child.Field(k).Name == "ID" {
				mIndex.id = k
			}
		}

		c := cron.New()
		var hasCron bool
		for k := 0; k < jobListRv.Len(); k++ {
			vv := jobListRv.Index(k)
			if vv.Kind() == reflect.Ptr {
				vv = vv.Elem()
			}
			id := vv.Field(mIndex.id).Uint()
			name := vv.Field(mIndex.name).String()
			express := vv.Field(mIndex.express).String()
			params := vv.Field(mIndex.parameters).String()
			key := vv.Field(mIndex.key).String()
			if express == "" {
				continue
			}
			if _, ok := jobs.Load(key); !ok {
				continue
			}
			f, _ := jobs.Load(key)

			v, _ := f.(*crontab)
			sched := &schedule{
				id:      id,
				jobName: name,
				params:  params,
				method:  v.Func,
			}

			err := c.AddJob(express, sched)
			if err != nil {
				logger.Error(err)
				continue
			}
			hasCron = true
		}
		if !hasCron {
			logger.Info("定时任务为空，跳过")
			return
		}
		logger.Info("定时任务开始执行...")
		c.Start()

		for {
			select {
			case <-ch:
				c.Stop()
				return
			}
		}
	}()
}

//Stop 停止定时任务
func Stop() {
	logger.Info("定时任务停止")
	if ch == nil {
		return
	}
	close(ch)

}

//Restart 重启
func Restart() {
	logger.Info("重启定时任务")
	close(ch)
	ch = nil
	Start()
}

//JobList 返回已注册的 定时任务列表
func JobList() (list map[string]string) {
	list = make(map[string]string)
	jobs.Range(func(k interface{}, v interface{}) bool {
		vv := v.(*crontab)
		list[k.(string)] = vv.Comment
		return true
	})
	return
}
