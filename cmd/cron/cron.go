package cron

import (
	"errors"
	"reflect"
	"sync"

	"github.com/bertang/bert/common/config/application"
	"github.com/bertang/bert/common/logger"
	"github.com/robfig/cron"
)

var (
	ch      chan interface{}  //用于停止定时任务的channel
	jobs    sync.Map          //注册任务的保存地
	service ICronService      //用于获取任务的service
	mIndex  = new(modelIndex) //用于反射执行时保存对应属性的位置
)

//modelIndex 用于保存反射属性的位置
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

//crontab 定时任务
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
	//如果是debug模式不开启定时任务
	if application.GetAppConf().Debug {
		return
	}
	//用于停止定时任务的channel
	ch = make(chan interface{})

	//执行定时任务的goroutine
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

		//判断是否是指针，如果是指针取其指向的实际数据
		if jobListType.Kind() == reflect.Ptr {
			jobListRv = jobListRv.Elem()
			jobListType = jobListType.Elem()
		}

		// 如果实际类型不是slice则退出
		if jobListType.Kind() != reflect.Slice {
			return
		}

		//判断slice保存的数据的实际类型，取期实际数据
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

			//哪果使用了 gorm.Model如法设置 tag所以使用属性名来获取
			if child.Field(k).Name == "ID" {
				mIndex.id = k
			}
		}

		//创建定时任务
		c := cron.New()
		var hasCron bool //用于判断是否有定时任务
		for k := 0; k < jobListRv.Len(); k++ {
			//从反射获取数据
			vv := jobListRv.Index(k)
			if vv.Kind() == reflect.Ptr {
				vv = vv.Elem()
			}

			//id必须是int类型
			id := vv.Field(mIndex.id).Uint()
			name := vv.Field(mIndex.name).String()
			express := vv.Field(mIndex.express).String()
			params := vv.Field(mIndex.parameters).String()
			key := vv.Field(mIndex.key).String()

			//如果时间表达式为空，则没有必要进行下去，跳过
			if express == "" {
				continue
			}

			//如果不存在这个key则跳过
			if f, ok := jobs.Load(key); ok {
				v, _ := f.(*crontab)
				err := c.AddJob(express, &schedule{
					id:      uint(id),
					jobName: name,
					params:  params,
					method:  v.Func,
				})
				if err != nil {
					logger.Error(err)
					continue
				}
				hasCron = true
			}
		}
		//如果没有定时任务 则跳出
		if !hasCron {
			logger.Info("定时任务为空，跳过")
			return
		}

		logger.Info("定时任务开始执行...")
		c.Start()

		//在另一个goroutine中 阻塞，
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
	ch = nil
}

//Restart 重启
func Restart() {
	logger.Info("重启定时任务")
	close(ch)
	//重置为空
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
