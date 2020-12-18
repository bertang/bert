package test

import (
	"fmt"
	"github.com/robfig/cron"
	"testing"
	"time"
)

func TestCronSchedule(t *testing.T) {
	c:= cron.New()
	//schedule:= &cron2.Schedule{
	//	ID:     1,
	//	Method: HasParam,
	//}
	//c.AddJob("*/5 * * * * ?", schedule)
	c.Start()

	for {
		time.Sleep(100 * time.Second)
	}
}


func NoParam() {
	fmt.Println("no Params...")
}

func HasParam(p ...interface{}) {
	fmt.Println(p...)
}