package logger

import (
        "io"
        "os"

        "github.com/bertang/bert/common/config/application"
        "github.com/kataras/golog"
        "github.com/natefinch/lumberjack"
)

var (
        lumberJackLogger *lumberjack.Logger
        logger           *golog.Logger
)

func init() {
        initLogger()
}
func initLogger() {
        appConf := application.GetAppConf()
        if appConf.Debug {
                logger = golog.New()
                logger.TimeFormat = "2006-01-02 15:04:05"
                logger.SetOutput(os.Stdout)
                return
        }
        lumberJackLogger = &lumberjack.Logger{
                Filename:   appConf.LoggerName,
                MaxSize:    appConf.MaxLogAge,
                MaxAge:     appConf.MaxLogAge,
                MaxBackups: appConf.MaxBackup,
                LocalTime:  true,
                Compress:   appConf.Compress,
        }
}

//SetLogger 设置记录日志
func SetLogger(gologger *golog.Logger) {
        logger = gologger
}

//GetWriter 获取日志相关格式
func GetWriter() io.Writer {
        appConf := application.GetAppConf()
        if lumberJackLogger == nil || appConf.Debug {
                return os.Stdout
        }
        return lumberJackLogger
}

//Info 信息
func Info(v ...interface{}) {
        logger.Info(v...)
}

//Warn 警告
func Warn(v ...interface{}) {
        logger.Warn(v...)
}

//Error 错误
func Error(v ...interface{}) {
        logger.Error(v...)
}

//Fatal 致命错误
func Fatal(v ...interface{}) {
        logger.Fatal(v...)
}

//Infof 格式化输出信息
func Infof(format string, v ...interface{}) {
        logger.Infof(format, v...)
}

//Fatalf 致命错误格式化
func Fatalf(format string, v ...interface{}) {
        logger.Fatalf(format, v...)
}

//Warnf 警告错误格式化
func Warnf(format string, v ...interface{}) {
        logger.Warnf(format, v...)
}

//Errorf 错误输出格式化
func Errorf(format string, v ...interface{}) {
        logger.Errorf(format, v...)
}