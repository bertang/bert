package mysql

import (
    "fmt"
    "log"

    "github.com/bertang/bert/common/config"
)

var (
    //multiConf 使用多数据库 默认的key为default
    multiConf  map[string]*mysqlConf
    dbFilePath = &config.FilePath{Filename: "mysql.yml"}
)
const (
    //DefaultDatabaseKey 默认数据库的key
    DefaultDatabaseKey = "default"
)

func initMysqlConf() {
    multiConf = make(map[string]*mysqlConf)
    config.Register(dbFilePath, &multiConf)
}

type mysqlConf struct {
    Host            string
    Port            int
    DBName          string `mapstructure:"db_name"`
    User            string
    Password        string
    Charset         string
    MaxConn         int    `mapstructure:"conns"`
    MaxIde          int    `mapstructure:"ides"`
    ConnMaxLifeTime int    `mapstructure:"max_lifetime"`
    SingularTable   bool   `mapstructure:"sigular"`
    TablePrefix     string `mapstructure:"prefix"`
}

//GetDatabaseKeys 获取数据库存储的key
func GetDatabaseKeys() []string {
    initMysqlConf()
    var keys []string
    for k := range multiConf {
        keys = append(keys, k)
    }
    return keys
}

//GetDSN 返回格式化的数据库dsn
func (m *mysqlConf) GetDSN() string {
    if m.Host == "" || m.User == "" || m.Password == "" {
        log.Fatal("数据库配置不正确！")
    }
    temp := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
    if m.Charset == "" {
        m.Charset = "utf8mb4"
    }
    return fmt.Sprintf(temp, m.User, m.Password, m.Host, m.Port, m.DBName, m.Charset)
}

//GetMysqlConf 获取mysql配置文件
func GetMysqlConf() *mysqlConf {
    if multiConf == nil || len(multiConf) == 0 {
        initMysqlConf()
    }
    if _, ok := multiConf[DefaultDatabaseKey]; ok {
        return multiConf[DefaultDatabaseKey]
    }
    return nil
}

//根据key获取配置文件
func GetMysqlConfByKey(key string) *mysqlConf {
    if multiConf == nil || len(multiConf) == 0 {
        initMysqlConf()
    }
    if _, ok := multiConf["key"]; ok {
        return multiConf["key"]
    }
    return nil
}

//SetMysqlConfFileName 自定义文件名
func SetMysqlConfFileName(name string) {
    dbFilePath.Filename = name
}
