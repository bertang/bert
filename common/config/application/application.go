package application

import (
    "path"

    "github.com/bertang/bert/common/config"
)

var (
    app     *application                            //获取到的app 缓存
    appPath = &config.FilePath{Filename: "app.yml"} //定义了位置
)

//application 一些app的配置
type application struct {
    Debug      bool     `mapstructure:"debug"`        //是否是开发模式
    UploadPath string   `mapstructure:"upload_path"`  //图片上传路径
    StaticDir  string   `mapstructure:"static_dir"`   //静态文件地址，主要是用于临时作为图片服务器的静态资源地址
    Host       string   `mapstructure:"host"`         //服务运行绑定ip
    Port       int      `mapstructure:"port"`         // 服务运行端口
    Name       string   `mapstructure:"name"`         //服务名称
    PageSize   int      `mapstructure:"page_size"`    //分页大小
    LoggerName string   `mapstructure:"logger_name"`  // 日志文件名称
    MaxLogSize int      `mapstructure:"max_log_size"` //单个日志文件最大大小 单位M
    MaxLogAge  int      `mapstructure:"max_log_age"`  //最长保存日志文件的时间 单位：天
    MaxBackup  int      `mapstructure:"max_backup"`   //最多可以保存多少个日志备份
    Compress   bool     `mapstructure:"log_compress"` //是否压缩日志
    JwtSecret  string   `mapstructure:"jwt_key"`      // jwt密码
    JwtMaxAge  int      `mapstructure:"jwt_max_age"`  //jwt最大的有效期,时间小时
    Cors       []string `mapstructure:"cors"`         //需要开启跨域的域名
    CorsMethod []string `mapstructure:"cors_method"`  //需要开启跨域的请求方法
    TimeFormat string   `mapstructure:"time_format"`  //时间格式化
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

    //日志相关
    if app.MaxLogAge == 0 {
        app.MaxLogAge = 100
    }
    if app.MaxBackup == 0 {
        app.MaxLogAge = 100
    }
    if app.MaxLogSize == 0 {
        app.MaxLogAge = 100
    }

    //默认的分页大小
    if app.PageSize == 0 {
        app.PageSize = 20
    }

    if app.UploadPath != "" && !path.IsAbs(app.UploadPath) {
        app.UploadPath = path.Join(config.GetAppPath(), app.UploadPath)
    }

    if app.TimeFormat == "" {
    	app.TimeFormat = "2006-01-02 15:04:05"
    }

}

//GetAppConf 获取系统配置
func GetAppConf() *application {
    if app == nil {
        initApplication()
    }
    return app
}

//SetApplicationConfigName 可以自定义的配置文件名
func SetApplicationConfigName(name string) {
    appPath.Filename = name
}
