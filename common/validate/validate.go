package validate

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/bertang/bert/common/logger"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func initValidate() (validate *validator.Validate, trans ut.Translator) {
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	//验证器注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		logger.Error("初始化中文验证器失败：%s", err)
	}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name != "" {
			return name
		}
		return fld.Name
	})
	_ = validate.RegisterValidation("mobile", CheckMobile)
	return
}

//ValidateStruct 验证结构体
//@data 需要验证的结构体
func Validate(data interface{}) map[string]string {
	//判断传入的类型是否是结构体
	rtData := reflect.TypeOf(data)
	if rtData.Kind() == reflect.Ptr {
		rtData = rtData.Elem()
	}

	if rtData.Kind() != reflect.Struct {
		panic("验证不是一个结构体")
	}
	validate, trans := initValidate()
	//验证
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	//有错误则对应一个字段与错误的关系
	fieldJSONTag := make(map[string]string, rtData.NumField())
	for k := 0; k < rtData.NumField(); k++ {
		if rtData.Field(k).Type.Kind() == reflect.Struct {
			rtSub := rtData.Field(k).Type
			for i := 0; i < rtSub.NumField(); i++ {
				fieldJSONTag[rtSub.Field(i).Name] = rtSub.Field(i).Tag.Get("json")
			}
		} else {
			fieldJSONTag[rtData.Field(k).Name] = rtData.Field(k).Tag.Get("json")
		}

	}
	errMap := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(err.StructField())
		errMap[fieldJSONTag[err.StructField()]] = err.Translate(trans)
	}
	return errMap
}

//CheckMobile 检查手机号
func CheckMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return  rgx.MatchString(fl.Field().String())
}
