//@Title redis.go
//@Description 
//@Author bertang
//@Created bertang 2021/2/7 11:52 上午
package redis

import (
    redis2 "github.com/bertang/bert/common/config/redis"
    "github.com/gomodule/redigo/redis"
    "sync"
)

var (
    rPool *redis.Pool
    once  sync.Once
)

//initRedis 初始化redis
func initRedis() {
    once.Do(func() {
        rConf := redis2.GetRedisConf()
        rPool = &redis.Pool{
            Dial: func() (redis.Conn, error) {
                c, err := redis.Dial("tcp", rConf.String())
                if err != nil {
                    return nil, err
                }
                if _, err := c.Do("auth", rConf.Auth); err != nil {
                    c.Close()
                    return nil, err
                }
                return c, err
            },
            MaxIdle:         rConf.MaxIdle,
            MaxActive:       rConf.MaxActive,
            IdleTimeout:     0,
            Wait:            false,
            MaxConnLifetime: 0,
        }
    })
}

//GetConn 获取redis连接
func GetConn() redis.Conn {
    if rPool == nil {
        initRedis()
    }
    return rPool.Get()
}
