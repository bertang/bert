package test

import (
	"github.com/kataras/golog"
	"go.uber.org/zap"
	"testing"
)

func BenchmarkGoLog(b *testing.B) {
	logger:= golog.New()
	for i:=0; i< b.N; i++ {
		logger.Info(i)
	}
}

func BenchmarkZapGoLog(b *testing.B) {
	loggerZap,_:= zap.NewProduction()
	loggerZap.Sugar()
}