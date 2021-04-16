package server

import (
        "fmt"
        "os"
        "path"

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
        } else {
                app = iris.New()
        }

        //服务常用配置
        app.Configure(iris.WithConfiguration(configuration()))
        //注册当服务停上时运行函数
        iris.RegisterOnInterrupt(onInterrupt)

        //设置跨域
        setCors()

        //添加静态文件支持
        if appConf.StaticDir != "" {
                app.HandleDir(path.Join("", appConf.StaticDir), iris.Dir(appConf.StaticDir))
        }

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

        //跨域
        c := cors.New(cors.Options{
                AllowedOrigins: appConf.Cors,
                AllowedMethods: appConf.CorsMethod,
                //需要添加所有的头
                AllowedHeaders:   []string{"Accept", "Origin", "XRequestedWith", "Content-Type", "LastModified", "Authorization"},
                AllowCredentials: true,
        })

        //跨域提前设置
        app.WrapRouter(c.ServeHTTP)

}

func onInterrupt() {
        cron.Stop()
        os.Exit(1)
}

func configuration() iris.Configuration {
        app:= application.GetAppConf()
        configuration := iris.Configuration{
                TimeFormat: app.TimeFormat,
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
