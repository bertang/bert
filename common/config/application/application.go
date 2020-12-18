package application

import (
	"github.com/bertang/bert/common/config"
)

var (
	app  *application

	appPath = &config.FilePath{Filename: "app.yml"}
)

type application struct {
	Debug      int    `mapstructure:"debug"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Name       string `mapstructure:"name"`
	LoggerName string `mapstructure:"logger_name"`
	MaxLogSize int    `mapstructure:"max_log_size"`
	MaxLogAge  int    `mapstructure:"max_log_age"`
	MaxBackup  int    `mapstructure:"max_backup"`
	Compress   bool   `mapstructure:"log_compress"`
	JwtSecret  string `mapstructure:"jwt_key"`
}


//初始化系统
func initApplication() {
	app = &application{}
	//注册配置
	config.Register(appPath, app)

	//初始化一些服务相关默认值
	if app.LoggerName == "" {
		app.LoggerName = "./logs/log.log"
	}
	if app.MaxLogAge == 0 {
		app.MaxLogAge = 50
	}
	if app.MaxBackup == 0 {
		app.MaxLogAge = 50
	}
	if app.MaxLogSize == 0 {
		app.MaxLogAge = 50
	}

}

//GetAppConf 获取系统配置
func GetAppConf() *application {
	if app == nil {
		initApplication()
	}
	return app
}

func SetApplicationConfigName(name string) {
	appPath.Filename = name
}