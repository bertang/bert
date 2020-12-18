package config

import (
	"fmt"
	"framework/common/files"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var (
	appPath string
	once    sync.Once
	confDirName = "conf"
)

type IConfig interface {
	getPath() string
}
type FilePath struct {
	Filename string
}

//初始化常用的路径等
func initConfPath() {
	fmt.Println("config...")
	once.Do(func() {
		var err error
		var workPath string

		if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			log.Fatal(err)
		}
		if workPath, err = os.Getwd(); err != nil {
			log.Fatal(err)
		}

		//兼容单元测试时
		if strings.HasSuffix(workPath, "test") {
			workPath = workPath[:strings.LastIndexByte(workPath, os.PathSeparator)]
		}

		//根据配置文件目录来获取项目路径，以获取正确的配置文件
		configPath := filepath.Join(appPath, confDirName)
		if !files.FileExists(configPath) {
			configPath = filepath.Join(workPath, confDirName)
			if !files.FileExists(configPath) {
				log.Fatalf("%s 配置文件目录不存在！", configPath)
			}
			appPath = workPath
		}
	})
}

func (c *FilePath) getPath() string {
	return c.Filename
}

//SetConfDirName 设置配置文件目录
func SetConfDirName(dirname string ) {
	confDirName = dirname
}

//Register 注册配置文件
func Register(config IConfig, conf interface{}) {
	if appPath == "" {
		initConfPath()
	}
	if conf == nil {
		log.Fatalln("非法的值")
	}
	configFilePath := filepath.Join(appPath, confDirName, config.getPath())
	fmt.Println(configFilePath)
	if !files.FileExists(configFilePath) {
		log.Fatalf("配置文件%s不存在！\n", configFilePath)
	}

	v := viper.New()
	v.SetConfigFile(configFilePath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalln("读取配置出错！", err)
	}

	//反射赋值
	rtConf := reflect.TypeOf(conf)
	if rtConf.Kind() != reflect.Ptr {
		log.Fatalln("非法的值")
	}
	//指针
	if rtConf.Elem().Name() == "Viper" {
		rvConf := reflect.ValueOf(conf)
		if !rvConf.Elem().CanAddr() {
			log.Fatalln("非法的值")
		}
		rvV := reflect.ValueOf(v)
		rvConf.Elem().Set(rvV.Elem())
		return
	}

	if err := v.Unmarshal(&conf); err != nil {
		log.Fatalln("解析配置文件出错")
	}
}
