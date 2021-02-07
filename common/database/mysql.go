package database

import (
	"log"

	mysqlConf "github.com/bertang/bert/common/config/mysql"
	"gorm.io/driver/mysql"

	"time"

	"github.com/bertang/bert/common/config/application"

	logger2 "github.com/bertang/bert/common/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	dbs map[string]*gorm.DB
)

func init() {
	databases := mysqlConf.GetDatabaseKeys()
	dbs = make(map[string]*gorm.DB)
	for _, v := range databases {
		initConnection(v)
	}
}

//初始化数据库连接
func initConnection(key string) {
	dbConf := mysqlConf.GetMysqlConfByKey(key)

	//配置
	conf := &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: dbConf.SingularTable, TablePrefix: dbConf.TablePrefix},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if application.GetAppConf().Debug {
		conf.Logger = logger.Default.LogMode(logger.Info)
	} else {
		conf.Logger = logger.Default.LogMode(logger.Error)
	}

	//连接数据库
	db, err := gorm.Open(mysql.Open(dbConf.GetDSN()), conf)
	if err != nil {
		logger2.Fatalf("初始化数据库失败", err)
		return
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

	dbs[key] = db
}

//GetDB 获取数据库连接
func GetDB() *gorm.DB {
	if _, ok := dbs[mysqlConf.DefaultDatabaseKey]; !ok {
		logger2.Errorf("没有当前数据库")
		return nil
	}
	return dbs[mysqlConf.DefaultDatabaseKey]
}

//GetDatabaseByKey 根据数据库的key获取数据
func GetDatabaseByKey(key string) *gorm.DB {
	if _, ok := dbs[key]; !ok {
		logger2.Errorf("没有当前数据库")
		return nil
	}
	return dbs[key]
}
