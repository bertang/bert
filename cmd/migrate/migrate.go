package migrate

import (
        "reflect"

        "github.com/bertang/bert/common/database"
        "github.com/bertang/bert/common/logger"
)

var (
        Migrates []*Migration
)

//Migration 迁移表
type Migration struct {
        Model      interface{}
        ColumnAdd  []string
        ColumnDrop []string
        Rename     [][]string
}

//Migrate 迁移数据表
func Migrate(models ...interface{}) {
        for _, v := range models {
                rt := reflect.TypeOf(v)
                if rt.Kind() == reflect.Ptr {
                        if rt.Elem().Kind() != reflect.Struct {
                                continue
                        }
                        if vv, ok := v.(*Migration); ok {
                                Migrates = append(Migrates, vv)
                        } else {
                                Migrates = append(Migrates, &Migration{Model: v})
                        }
                } else if rt.Kind() == reflect.Struct {
                        if vv, ok := v.(Migration); ok {
                                Migrates = append(Migrates, &vv)
                        } else {
                                Migrates = append(Migrates, &Migration{Model: &vv})
                        }
                }

        }
}

//Start 执行数据库迁移
func Start() {
        logger.Info("开始迁移数据库。。。")
        migrator := database.GetDB().Migrator()
        for _, v := range Migrates {

                if !migrator.HasTable(v.Model) {
                        _ = migrator.AutoMigrate(v.Model)
                }

                //添加字段
                for _, c := range v.ColumnAdd {
                        if !migrator.HasColumn(v.Model, c) {
                                _ = migrator.AddColumn(v.Model, c)
                        }
                }

                //删除字段
                for _, c := range v.ColumnDrop {
                        if migrator.HasColumn(v.Model, c) {
                                _ = migrator.DropColumn(v.Model, c)
                        }

                }

                //重命名字段
                for _, c := range v.Rename {
                        if migrator.HasColumn(v.Model, c[0]) && !migrator.HasColumn(v.Model, c[1]) {
                                _ = migrator.RenameColumn(v.Model, c[0], c[1])
                        }
                }

        }

        logger.Info("数据库迁移结束...")
}
