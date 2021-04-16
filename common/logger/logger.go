package logger

import (
    "os"
    "runtime"

    "github.com/bertang/bert/common/config/application"
    "github.com/kataras/golog"
    "github.com/natefinch/lumberjack"
)

var (
    logger           = golog.Default
)

func init() {
    appConf := application.GetAppConf()

    logger.SetTimeFormat(appConf.TimeFormat)
    if appConf.Debug {
        logger.SetOutput(os.Stdout)
        return
    }
    writer := &lumberjack.Logger{
        Filename:   appConf.LoggerName,
        MaxSize:    appConf.MaxLogAge,
        MaxAge:     appConf.MaxLogAge,
        MaxBackups: appConf.MaxBackup,
        LocalTime:  true,
        Compress:   appConf.Compress,
    }
    logger.SetOutput(writer)

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

//LoggerErrInfo 记录错误位置
func ErrInfo(v ...interface{}) {
    pc, _,_,_ := runtime.Caller(1)
    f:=runtime.FuncForPC(pc)
    Errorf("发生错误：" + f.Name() + ":%s", v...)
}