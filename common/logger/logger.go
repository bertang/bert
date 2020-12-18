package logger

import (
	"framework/common/config/application"
	"github.com/kataras/golog"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
)

var (
	lumberJackLogger *lumberjack.Logger
	logger           *golog.Logger
)

func initLogger() {
	appConf := application.GetAppConf()
	if appConf.Debug == 1 {
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

//设置记录日志
func SetLogger(gologger *golog.Logger) {
	logger = gologger
}

//获取日志相关格式
func GetLogger() io.Writer {
	appConf := application.GetAppConf()
	if lumberJackLogger == nil || appConf.Debug == 1 {
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
