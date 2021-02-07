//@Title redis.go
//@Description 
//@Author bertang
//@Created bertang 2021/2/7 11:53 上午
package redis

import (
    "fmt"
    "github.com/bertang/bert/common/config"
    "sync"
)

var (
    redisConf  *RedisConf
    once       sync.Once
    dbFilePath = &config.FilePath{Filename: "redis.yml"}
)

//initRedis 初始化redis
func initRedis() {
    once.Do(func() {
        redisConf = &RedisConf{}
        config.Register(dbFilePath, redisConf)
    })

}

//RedisConf redis配置
type RedisConf struct {
    Host      string
    Port      int
    Auth      string
    MaxIdle   int `mapstructure:"max_idle"`
    MaxActive int `mapstructure:"max_active"`
}

//String 返回
func (r *RedisConf) String() string {
    return fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
}

func GetRedisConf() *RedisConf {
    if redisConf == nil {
        initRedis()
    }
    return redisConf
}
