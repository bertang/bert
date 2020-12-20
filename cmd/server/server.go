package server

import (
        "fmt"
        "github.com/bertang/bert/cmd/cron"
        "github.com/bertang/bert/cmd/server/router"
        "github.com/bertang/bert/common/config/application"
        "github.com/bertang/bert/common/logger"

        "github.com/kataras/iris/v12"
        "github.com/rs/cors"
)

var (
        app     *iris.Application
        appConf = application.GetAppConf()
)

//Exec web运行
func Exec() {
        run()
}

func run() {

        if appConf.Debug {
                app = iris.Default()
        } else {
                app = iris.New()
        }

        //服务常用配置
        app.Configure(iris.WithConfiguration(configuration()))
        //注册当服务停上时运行函数
        iris.RegisterOnInterrupt(onInterrupt)

        //设置日志
        app.Logger().SetOutput(logger.GetWriter())

        //设置跨域
        setCors()

        //路由处理
        router.RegisterRouter(app)

        //错误处理
        app.OnAnyErrorCode(onCodeErrorHandle)

        //运行服务
        err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", appConf.Host, appConf.Port)))
        if err != nil {
                logger.Fatal(err)
        }
}

//错误处理
func onCodeErrorHandle(ctx iris.Context) {
        //_, err := ctx.JSON(tools.Failure(ctx.GetStatusCode()))
        //logger.Error(err)
        ctx.StopExecution()
}

//设置跨域
func setCors() {
        c := cors.New(cors.Options{
                AllowedOrigins:   nil,
                AllowedMethods:   nil,
                AllowedHeaders:   nil,
                AllowCredentials: false,
        })
        //跨域提前设置
        app.WrapRouter(c.ServeHTTP)
}

func onInterrupt() {
        cron.Stop()
}

func configuration() iris.Configuration {
        configuration := iris.Configuration{
                TimeFormat: "2006-01-02 15:04:05",
                Charset:    "UTF-8",
        }

        if appConf.Debug{
                configuration.EnableOptimizations = true
                configuration.FireMethodNotAllowed = true
                configuration.DisableStartupLog = true
                configuration.IgnoreServerErrors = []string{iris.ErrServerClosed.Error()}
        }

        return configuration
}
