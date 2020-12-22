package server

import (
        "fmt"
        "os"

        "github.com/bertang/bert/cmd/cron"
        "github.com/bertang/bert/cmd/server/router"
        "github.com/bertang/bert/common/config/application"
        "github.com/bertang/bert/common/logger"
        "github.com/kataras/iris/v12/context"
        "github.com/rs/cors"

        "github.com/kataras/iris/v12"
)

var (
        app           *iris.Application
        appConf       = application.GetAppConf()
        errCodeHander []context.Handler
)

//Start web运行
func Start() {
        run()
}

//SetOnErrCodeHandler 设置错误处理
func SetOnErrCodeHandler(handler ...context.Handler) {
        errCodeHander = handler
}

func run() {
        if appConf.Debug {
                app = iris.Default()
        } else {
                app = iris.New()
                app.Logger().SetLevel("error")
        }

        //服务常用配置
        app.Configure(iris.WithConfiguration(configuration()))
        //注册当服务停上时运行函数
        iris.RegisterOnInterrupt(onInterrupt)

        //设置日志
        app.Logger().SetOutput(logger.GetWriter())
        //统一为一个logger
        logger.SetLogger(app.Logger())

        //设置跨域
        setCors()

        //路由处理
        router.RegisterRouter(app)

        //错误处理
        if len(errCodeHander) > 0 {
                app.OnAnyErrorCode(errCodeHander...)
        }

        //运行服务
        err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", appConf.Host, appConf.Port)))
        if err != nil {
                logger.Fatal(err)
        }
}

//设置跨域
func setCors() {
        if len(appConf.Cors) == 0 {
                return
        }
        c := cors.New(cors.Options{
                AllowedOrigins:   appConf.Cors,
                AllowedMethods:   appConf.CorsMethod,
                AllowedHeaders:   nil,
                AllowCredentials: false,
        })
        //跨域提前设置
        app.WrapRouter(c.ServeHTTP)

}

func onInterrupt() {
        cron.Stop()
        os.Exit(1)
}

func configuration() iris.Configuration {
        configuration := iris.Configuration{
                TimeFormat: "2006-01-02 15:04:05",
                Charset:    "UTF-8",
        }

        if appConf.Debug {
                configuration.EnableOptimizations = true
                configuration.FireMethodNotAllowed = true
                configuration.DisableStartupLog = true
                configuration.IgnoreServerErrors = []string{iris.ErrServerClosed.Error()}
        }

        return configuration
}
