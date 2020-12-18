package env

import (
	"framework/common/config"
	"github.com/spf13/viper"
)

var (
	env *viper.Viper
)

//初始化自定义的配置
func initEnv() {
	configPath := &config.FilePath{Filename: "env.yml"}
	env := viper.New()
	config.Register(configPath, &env)
}

//Env 获取string 类型的配置
func Env(key string) string {
	if env == nil {
		initEnv()
	}
	return env.GetString(key)
}

//Viper 获取viper原始对象
func Viper() *viper.Viper {
	return env
}

//EnvInt 获取int型配置
func EnvInt(key string) int {
	if env == nil {
		initEnv()
	}
	return env.GetInt(key)
}

//EnvBool 获取bool型
func EnvBool(key string) bool {
	return viper.GetBool(key)
}
