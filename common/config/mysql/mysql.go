package mysql

import (
	"fmt"
	"log"

	"github.com/bertang/bert/common/config"
)

var (
	mConf      *mysqlConf
	dbFilePath = &config.FilePath{Filename: "mysql.yml"}
)

func initMysqlConf() {
	mConf = new(mysqlConf)
	config.Register(dbFilePath, &mConf)
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
	initMysqlConf()
	return mConf
}

//SetMysqlConfFileName 自定义文件名
func SetMysqlConfFileName(name string) {
	dbFilePath.Filename = name
}
