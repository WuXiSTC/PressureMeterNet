package option

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
)

type Option struct {
	ServerInfo          ServerInfoOption    `yaml:"ServerInfo" usage:"Necessary information about this server."`
	GogisnetOption      GogisnetOption      `yaml:"GogisnetOption" usage:"Option for Gogisnet."`
	GRPCListenerOption  ListenerOption      `yaml:"ListenerOption" usage:"Option for port listen when running the server."`
	PressureMeterConfig PressureMeterConfig `yaml:"PressureMeterConfig" usage:"Configuration for PressureMeter server."`
}

func GenerateOption(logger func(i ...interface{})) (opt Option) {
	opt = Option{
		ServerInfo:          DefaultServerInfoOption(),
		GogisnetOption:      DefaultGogisnetOption(),
		GRPCListenerOption:  DefaultListenerOption(),
		PressureMeterConfig: DefaultPressureMeterConfig(),
	}
	file := flag.String("c", "", "Set this value to load option from a option file at this path (Must be a .yaml file).")
	gfile := flag.String("g", "", "Set this value to generate a option file.")
	generateStructOption(opt, "", "")
	flag.Parse()
	if path := *file; path != "" {
		logger(fmt.Sprintf("-c flag detected, load option from %s.", path)) //文件读取模式
		data, err := ioutil.ReadFile(path)                                  //读取文件
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &opt) //反序列化
		if err != nil {
			panic(err)
		}
	}
	if path := *gfile; path != "" {
		logger(fmt.Sprintf("-g flag detected, option file will write to %s.", path))
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
	}
	printStructOption(opt, "", func(i ...interface{}) { logger(i...) })
	return
}

//通过反射获取设置信息
func generateStructOption(option interface{}, OptionName, OptionComment string) {
	operateStruct(option, OptionName, OptionComment,
		func(Struct interface{}, StructName, StructUsage string) {
			StructValue := reflect.ValueOf(Struct)
			StructKind := StructValue.Elem().Kind()
			switch StructKind {
			case reflect.Uint64:
				pointer := StructValue.Interface().(*uint64)
				flag.Uint64Var(pointer, StructName, *pointer, StructUsage)
			case reflect.String:
				pointer := StructValue.Interface().(*string)
				flag.StringVar(pointer, StructName, *pointer, StructUsage)
			case reflect.Float64:
				pointer := StructValue.Interface().(*float64)
				flag.Float64Var(pointer, StructName, *pointer, StructUsage)
			}
		})
}

func printStructOption(option interface{}, OptionName string, logger func(...interface{})) {
	operateStruct(option, OptionName, "",
		func(Struct interface{}, StructName, StructUsage string) {
			StructValue := reflect.ValueOf(Struct).Elem()
			logger(StructName, "=", StructValue.Interface(), "\n")
		})
}

func operateStruct(Struct interface{}, StructName, StructUsage string,
	operation func(interface{}, string, string)) {
	StructType := reflect.TypeOf(Struct)
	StructValue := reflect.ValueOf(Struct)
	for i := 0; i < StructType.NumField(); i++ { // 遍历结构体所有成员
		field := StructType.Field(i)       //获取每个成员
		fieldType := field.Type            // 获取每个成员的结构体字段信息
		fieldValue := StructValue.Field(i) //获取每个成员的结构体字段值
		fieldKind := fieldType.Kind()      ////获取每个成员的结构体字段Kind

		fieldName := field.Name //获取字段名称
		if StructName != "" {
			fieldName = StructName + "." + field.Name
		}
		fieldTag := field.Tag.Get("usage")
		if StructUsage != "" {
			fieldTag = StructUsage + " " + fieldTag
		}

		if fieldKind == reflect.Struct { //如果还是一个结构体
			operateStruct(fieldValue.Interface(), fieldName, fieldTag, operation) //就递归
		} else { //不是结构体，那么执行操作
			operation(fieldValue.Interface(), fieldName, fieldTag)
		}
	}
}
