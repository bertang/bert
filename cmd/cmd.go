package cmd

import (
	"github.com/bertang/bert/cmd/cron"
	"github.com/bertang/bert/cmd/migrate"
	"github.com/bertang/bert/cmd/server"
)

type command func()

var (
	funs = []command{}
)

func init() {
	//先迁移，后服务
	addCommand(migrate.Start)
	// 定时任务启动
	addCommand(cron.Start)
	//api服务启动
	addCommand(server.Start)

}

func addCommand(cmd command) {
	funs = append(funs, cmd)
}

//Start 开始迁移
func Start() {
	for _, v := range funs {
		v()
	}
}
