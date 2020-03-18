package option

import (
	"flag"
	"fmt"
	"reflect"
)

func GenerateOption() (opt Option, lopt ListenerOption) {
	opt, lopt = DefaultOption(), DefaultListenerOption()
	generateStructOption(opt, "Option", "Comment")
	generateStructOption(lopt, "ListenerOption", "Comment")
	flag.Parse()
	printStructOption(opt, "Option")
	printStructOption(lopt, "ListenerOption")
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

func printStructOption(option interface{}, OptionName string) {
	operateStruct(option, OptionName, "",
		func(Struct interface{}, StructName, StructUsage string) {
			StructValue := reflect.ValueOf(Struct).Elem()
			fmt.Print(StructName, "=", StructValue.Interface(), "\n")
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

		fieldName := StructName + "." + field.Name //获取字段名称
		fieldTag := StructUsage + "-" + field.Tag.Get("usage")

		if fieldKind == reflect.Struct { //如果还是一个结构体
			operateStruct(fieldValue.Interface(), fieldName, fieldTag, operation) //就递归
		} else { //不是结构体，那么执行操作
			operation(fieldValue.Interface(), fieldName, fieldTag)
		}
	}
}
