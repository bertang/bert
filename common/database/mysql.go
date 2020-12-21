package database

import (
        mysqlConf "github.com/bertang/bert/common/config/mysql"
        "gorm.io/driver/mysql"
        "log"

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

        db, err = gorm.Open(mysql.Open(dbConf.GetDSN()), &gorm.Config{})
        if err != nil {
                log.Fatal(err)
        }
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

        db.NamingStrategy = schema.NamingStrategy{SingularTable: dbConf.SingularTable, TablePrefix: dbConf.TablePrefix}

        //配置是否是开发模式
        if application.GetAppConf().Debug {
                db.Logger = logger.Default.LogMode(logger.Info)
        }
}

//获取数据库连接
func GetDB() *gorm.DB {
        if db == nil {
                initConnection()
        }
        return db
}
