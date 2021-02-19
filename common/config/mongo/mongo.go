//@Title mongo.go
//@Description 
//@Author bertang
//@Created bertang 2021/2/17 6:45 下午
package mongo

import (
    "fmt"
    "github.com/bertang/bert/common/config"
    "sync"
)

var (
    mongoConf  *Conf
    once       sync.Once
    dbFilePath = &config.FilePath{Filename: "mongo.yml"}
)

//initRedis 初始化redis
func initMongoConf() {
    once.Do(func() {
        mongoConf = &Conf{}
        config.Register(dbFilePath, mongoConf)
    })

}

//MongoConf redis配置
type Conf struct {
    Host     string
    Port     int
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    Timeout  int    `mapstructure:"timeout"`
}

func (c Conf) String() string {
    return fmt.Sprintf("mongodb://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
}

//GetRedisConf 获取redis配置对象
func GetMongoConf() *Conf {
    if mongoConf == nil {
        initMongoConf()
    }
    return mongoConf
}
