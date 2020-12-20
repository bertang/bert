package migrate

import (
        "github.com/bertang/bert/common/db"
        "github.com/bertang/bert/common/logger"
        "reflect"
)

//迁移表
type migration struct {
        Model      interface{}
        ColumnAdd  []string
        ColumnDrop []string
        Rename     [][]string
}

//Exec 执行数据库迁移
func Exec() {
        logger.Info("开始迁移数据库。。。")
        db := db.GetDB()
        //TODO 获取 需要迁移的数据
        m:= make([]interface{}, 100)
        for _, v := range m {
                rt := reflect.TypeOf(v)
                if rt.Kind() != reflect.Ptr {
                        logger.Fatal("只接受指针类型")
                }
                if rt.Elem().Name() == "migration" {
                        vv := v.(*migration)
                        if !db.Migrator().HasTable(vv.Model) {
                                _ = db.AutoMigrate(vv.Model)
                        }

                        //添加字段
                        for _, c := range vv.ColumnAdd {
                                if !db.Migrator().HasColumn(vv.Model, c) {
                                        _ = db.Migrator().AddColumn(vv.Model, c)
                                }
                        }

                        //删除字段
                        for _, c := range vv.ColumnDrop {
                                if db.Migrator().HasColumn(vv.Model, c) {
                                        _ = db.Migrator().DropColumn(vv.Model, c)
                                }

                        }

                        //重命名字段
                        for _, c := range vv.Rename {
                                if db.Migrator().HasColumn(vv.Model, c[0]) && !db.Migrator().HasColumn(vv.Model, c[1]) {
                                        _ = db.Migrator().RenameColumn(vv.Model, c[0], c[1])
                                }
                        }

                } else if rt.Elem().PkgPath() == "manager/datamodels" {
                        //只迁移datamodel包中的数据
                        if !db.Migrator().HasTable(v) {
                                db.AutoMigrate(v)
                        }

                } else {
                        logger.Fatalf("不支持的类型%s\n", rt.Elem().PkgPath())
                }

        }
        logger.Info("数据库迁移结束...")
}
