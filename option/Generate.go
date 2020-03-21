package option

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
)

func Generate(logger func(i ...interface{})) (opt Option, exit bool) {
	opt, exit = DefaultOption(), false
	file := flag.String("UseConfigFile", "", "Set this value to load option from a config file at this path (Must be a .yaml file).")
	gfile := flag.String("GenerateConfigFile", "", "Set this value to generate a option file.")
	generateValue(reflect.ValueOf(&opt), "", "")
	flag.Parse()
	if path := *file; path != "" {
		logger(fmt.Sprintf("-UseConfigFile flag detected, load option from %s.", path)) //文件读取模式
		data, err := ioutil.ReadFile(path)                                              //读取文件
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &opt) //反序列化
		if err != nil {
			panic(err)
		}
	}
	if path := *gfile; path != "" {
		logger(fmt.Sprintf("-GenerateConfigFile flag detected, option file will write to %s.", path))
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm) //打开文件流
		if err != nil {
			panic(err)
		}
		data, err := yaml.Marshal(opt) //序列化
		if err != nil {
			panic(err)
		}
		n, err := f.Write(data) //写入
		if err != nil {
			panic(err)
		}
		logger(fmt.Sprintf("%d bytes written to %s", n, path))
		exit = true
		return
	}
	printValue(reflect.ValueOf(&opt), "", func(i ...interface{}) { logger(i...) })
	return
}

//通过反射获取设置信息
func generateValue(Value reflect.Value, ValueName, ValueUsage string) {
	operateValue(Value, ValueName, ValueUsage,
		func(Value reflect.Value, ValueName, ValueUsage string) {
			switch Value.Kind() {
			case reflect.Uint64:
				pointer := Value.Addr().Interface().(*uint64)
				flag.Uint64Var(pointer, ValueName, *pointer, ValueUsage)
			case reflect.String:
				pointer := Value.Addr().Interface().(*string)
				flag.StringVar(pointer, ValueName, *pointer, ValueUsage)
			case reflect.Float64:
				pointer := Value.Addr().Interface().(*float64)
				flag.Float64Var(pointer, ValueName, *pointer, ValueUsage)
			}
		})
}

func printValue(Value reflect.Value, ValueName string, logger func(...interface{})) {
	operateValue(Value, ValueName, "",
		func(Value reflect.Value, ValueName, ValueUsage string) {
			p := func(TypeStr string) {
				logger(ValueName, fmt.Sprintf("=(%s)", TypeStr), Value.Interface())
			}
			switch Value.Type().Kind() {
			case reflect.String:
				p("string")
			case reflect.Float64:
				p("float64")
			case reflect.Uint64:
				p("uint64")
			case reflect.Int64:
				p("int64")
			case reflect.Bool:
				p("bool")
			}
		})
}

func operateValue(Value reflect.Value, ValueName, ValueUsage string,
	operation func(Value reflect.Value, ValueName, ValueUsage string)) {
	Value = Value.Elem()
	Type := Value.Type()
	for i := 0; i < Value.NumField(); i++ { // 遍历结构体所有成员
		field := Type.Field(i) // 获取每个成员的结构体字段

		fieldName := field.Name //获取字段名称
		if ValueName != "" {
			fieldName = ValueName + "." + field.Name
		}
		fieldTag := field.Tag.Get("usage")
		if ValueUsage != "" {
			fieldTag = ValueUsage + " " + fieldTag
		}

		fieldValue := Value.Field(i)             //获取每个成员的结构体字段值
		if field.Type.Kind() == reflect.Struct { //如果还是一个结构体
			operateValue(fieldValue.Addr(), fieldName, fieldTag, operation) //就递归
		} else { //不是结构体，那么执行操作
			operation(fieldValue, fieldName, fieldTag)
		}
	}
}
