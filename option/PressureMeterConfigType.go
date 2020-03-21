package option

import (
	"gitee.com/WuXiSTC/PressureMeter"
	"gitee.com/WuXiSTC/PressureMeter/Model"
	"reflect"
	"regexp"
	"strings"
)

type PressureMeterConfig struct {
	ModelConfig Model.Config `yaml:"ModelConfig" usage:"Configuration for Model in PressureMeter."`
	URLConfig   URLConfig    `yaml:"URLConfig" usage:"Configuration for URL of iris server in PressureMeter."`
}

func DefaultPressureMeterConfig() PressureMeterConfig {
	return PressureMeterConfig{
		ModelConfig: Model.DefaultConfig(),
		URLConfig:   defaultURLConfig(),
	}
}
func (o PressureMeterConfig) PutConfig(op *PressureMeter.Config) {
	op.ModelConfig = o.ModelConfig
	o.URLConfig.PutOption(&op.URLConfig)
}

type URLConfig struct {
	NewTask    string `yaml:"NewTask" usage:"URL for creating new task."`
	DeleteTask string `yaml:"DeleteTask" usage:"URL for deleting task."`
	GetConfig  string `yaml:"GetConfig" usage:"URL for getting jmx file of the task."`
	GetResult  string `yaml:"GetResult" usage:"URL for getting jtl file of the task."`
	GetLog     string `yaml:"GetLog" usage:"URL of for getting log file of the task."`
	StartTask  string `yaml:"StartTask" usage:"URL for starting a task."`
	StopTask   string `yaml:"StopTask" usage:"URL for stopping a task."`
	GetState   string `yaml:"GetState" usage:"URL for getting the running state of a task."`
}

func defaultURLConfig() URLConfig {
	URLConfig := URLConfig{}
	ConfigValue := reflect.ValueOf(&URLConfig).Elem()

	PMURLConfig := PressureMeter.DefaultConfig().URLConfig
	PMConfigValue := reflect.ValueOf(PMURLConfig)
	PMConfigType := reflect.TypeOf(PMURLConfig)

	for i := 0; i < PMConfigValue.NumField(); i++ {
		PMFieldValue := strings.Join(PMConfigValue.Field(i).Interface().([]string), "/")
		PMFieldName := PMConfigType.Field(i).Name
		ConfigValue.FieldByName(PMFieldName).Set(reflect.ValueOf(PMFieldValue))
	}
	return URLConfig
}
func (o URLConfig) PutOption(op *PressureMeter.URLConfig) {
	ConfigValue := reflect.ValueOf(o)

	PMConfigValue := reflect.ValueOf(op).Elem()
	PMConfigType := reflect.TypeOf(*op)

	for i := 0; i < PMConfigValue.NumField(); i++ {
		PMFieldName := PMConfigType.Field(i).Name
		FieldValue := ConfigValue.FieldByName(PMFieldName)
		PMFieldValueStr := FieldValue.Interface().(string)
		PMFieldValueStr = re.ReplaceAllString(PMFieldValueStr, "")    //去除非法字符
		PMFieldValueStr = rex.ReplaceAllString(PMFieldValueStr, "/")  //去除多余斜杠
		PMFieldValueStr = rebeg.ReplaceAllString(PMFieldValueStr, "") //去除开头斜杠
		PMFieldValueStr = reend.ReplaceAllString(PMFieldValueStr, "") //去除结尾斜杠
		PMFieldValue := strings.Split(PMFieldValueStr, "/")
		PMConfigValue.Field(i).Set(reflect.ValueOf(PMFieldValue))
	}
}

var re, _ = regexp.Compile("[^0-9a-zA-Z/]*")
var rex, _ = regexp.Compile("/+")
var rebeg, _ = regexp.Compile("^/*")
var reend, _ = regexp.Compile("/*$")
