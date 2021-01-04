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
        app            *iris.Application
        appConf        = application.GetAppConf()
        errCodeHandler []context.Handler
)

//Start web运行
func Start() {
        run()
}

//RegisterErrCodeHandler 设置错误处理
//@handler 执行的中间件
func RegisterErrCodeHandler(handler ...context.Handler) {
        errCodeHandler = handler
}

//执行
func run() {
        if appConf.Debug {
                app = iris.Default()
                app.Logger().SetOutput(os.Stdin)
        } else {
                app = iris.New()
                app.Logger().SetLevel("error")
        }

        //服务常用配置
        app.Configure(iris.WithConfiguration(configuration()))
        //注册当服务停上时运行函数
        iris.RegisterOnInterrupt(onInterrupt)

        //设置日志
        logger.SetLogger(app.Logger())

        //设置跨域
        setCors()

        //路由处理
        router.RegisterRouter(app)

        //错误处理
        if len(errCodeHandler) > 0 {
                app.OnAnyErrorCode(errCodeHandler...)
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

        //如果是生产环境，添加一些优化
        if appConf.Debug {
                configuration.EnableOptimizations = true
                configuration.FireMethodNotAllowed = true
                configuration.DisableStartupLog = true
                configuration.IgnoreServerErrors = []string{iris.ErrServerClosed.Error()}
        }

        return configuration
}
