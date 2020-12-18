package test

import (
	"fmt"
	"framework/common/config"
	"framework/common/config/application"
	"framework/common/config/mysql"
	"testing"
)

func TestApplication(t *testing.T) {
	application.SetApplicationConfigName("fuck.yml")
	config.SetConfDirName("config")
	app:= application.GetAppConf()
	fmt.Println(app)
}


func TestMysqlConf(t *testing.T) {
	config.SetConfDirName("configs")
	dbConf:= mysql.GetMysqlConf()
	fmt.Println(dbConf)
}