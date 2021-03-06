package migrate

import (
	"os"
	"reflect"
	"strings"

	"github.com/bertang/bert/common/config"
	"github.com/bertang/bert/common/config/application"

	"github.com/bertang/bert/common/database"
	"github.com/bertang/bert/common/logger"
)

var (
	//Migrates 暂存迁移数据
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

	//搜索在根目录下的sql文件
	if application.GetAppConf().Debug {
		writeSql()
	}
	logger.Info("数据库迁移结束...")
}

//搜索在根目录下的sql文件
func writeSql() {
	dirs, _ := os.ReadDir(config.GetAppPath())
	db := database.GetDB()
	for _, fileInfo := range dirs {
		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".sql") {
			continue
		}
		filePath := config.GetAppPath() + string(os.PathSeparator) + fileInfo.Name()
		sqlBytes, _ := os.ReadFile(filePath)
		sql := strings.Split(string(sqlBytes), ";")
		for _, v := range sql {
			db.Exec(v)
		}
	}
}
