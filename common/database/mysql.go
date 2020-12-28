package database

import (
	"log"

	mysqlConf "github.com/bertang/bert/common/config/mysql"
	"gorm.io/driver/mysql"

	"time"

	"github.com/bertang/bert/common/config/application"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
)

func init() {
	initConnection()
}

//初始化数据库连接
func initConnection() {
	var err error
	dbConf := mysqlConf.GetMysqlConf()

	//配置
	conf := &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: dbConf.SingularTable, TablePrefix: dbConf.TablePrefix},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if application.GetAppConf().Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	}

	//连接数据库
	db, err = gorm.Open(mysql.Open(dbConf.GetDSN()), conf)
	if err != nil {
		log.Fatal(err)
	}

	//连接池设置
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	if dbConf.MaxConn > 0 {
		sqlDb.SetMaxOpenConns(dbConf.MaxConn)
	}
	if dbConf.MaxIde > 0 {
		sqlDb.SetMaxIdleConns(dbConf.MaxIde)
	}
	if dbConf.ConnMaxLifeTime > 0 {
		sqlDb.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime) * time.Minute)
	}
}

//获取数据库连接
func GetDB() *gorm.DB {
	if db == nil {
		initConnection()
	}
	return db
}
