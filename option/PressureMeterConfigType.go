package option

import (
	"gitee.com/WuXiSTC/PressureMeter"
	"gitee.com/WuXiSTC/PressureMeter/Model"
	"gitee.com/WuXiSTC/PressureMeter/Model/Daemon"
	"gitee.com/WuXiSTC/PressureMeter/Model/Task"
	"reflect"
	"regexp"
	"strings"
)

type PressureMeterConfig struct {
	ModelConfig ModelConfig
	URLConfig   URLConfig
}

func DefaultPressureMeterConfig() PressureMeterConfig {
	return PressureMeterConfig{
		ModelConfig: DefaultModelConfig(),
		URLConfig:   DefaultURLConfig(),
	}
}
func (o PressureMeterConfig) PutConfig(op *PressureMeter.Config) {
	o.ModelConfig.PutOption(&op.ModelConfig)
	o.URLConfig.PutOption(&op.URLConfig)
}

type ModelConfig struct {
	DaemonConfig DaemonConfig
	TaskConfig   TaskConfig
}

func DefaultModelConfig() ModelConfig {
	return ModelConfig{
		DaemonConfig: DefaultDaemonConfig(),
		TaskConfig:   DefaultTaskConfig(),
	}
}
func (o ModelConfig) PutOption(op *Model.Config) {
	o.DaemonConfig.PutOption(&op.DaemonConfig)
	o.TaskConfig.PutOption(&op.TaskConfig)
}

type DaemonConfig struct {
	TaskAccN  *uint64
	TaskQSize *uint64
}

func DefaultDaemonConfig() DaemonConfig {
	TaskAccN := uint64(4)
	TaskQSize := uint64(16)
	return DaemonConfig{
		TaskAccN:  &TaskAccN,
		TaskQSize: &TaskQSize,
	}
}
func (o DaemonConfig) PutOption(op *Daemon.Config) {
	op.TaskAccN = *o.TaskAccN
	op.TaskQSize = *o.TaskQSize
}

type TaskConfig struct {
	JmxDir *string
	JtlDir *string
	LogDir *string
}

func DefaultTaskConfig() TaskConfig {
	JmxDir, JtlDir, LogDir := "Data/jmx", "Data/jtl", "Data/log"
	return TaskConfig{
		JmxDir: &JmxDir,
		JtlDir: &JtlDir,
		LogDir: &LogDir,
	}
}
func (o TaskConfig) PutOption(op *Task.Config) {
	op.JmxDir = *o.JmxDir
	op.JtlDir = *o.JtlDir
	op.LogDir = *o.LogDir
}

type URLConfig struct {
	NewTask    *string
	DeleteTask *string
	GetConfig  *string
	GetResult  *string
	GetLog     *string
	StartTask  *string
	StopTask   *string
	GetState   *string
}

func DefaultURLConfig() URLConfig {
	URLConfig := URLConfig{
		NewTask:    new(string),
		DeleteTask: new(string),
		GetConfig:  new(string),
		GetResult:  new(string),
		GetLog:     new(string),
		StartTask:  new(string),
		StopTask:   new(string),
		GetState:   new(string),
	}
	ConfigValue := reflect.ValueOf(URLConfig)

	PMURLConfig := PressureMeter.DefaultConfig().URLConfig
	PMConfigValue := reflect.ValueOf(PMURLConfig)
	PMConfigType := reflect.TypeOf(PMURLConfig)

	for i := 0; i < PMConfigValue.NumField(); i++ {
		PMFieldValue := strings.Join(PMConfigValue.Field(i).Interface().([]string), "/")
		PMFieldName := PMConfigType.Field(i).Name
		*(ConfigValue.FieldByName(PMFieldName).Interface().(*string)) = PMFieldValue
	}
	return URLConfig
}
func (o URLConfig) PutOption(op *PressureMeter.URLConfig) {
	ConfigValue := reflect.ValueOf(o)

	PMConfigValue := reflect.ValueOf(*op)
	PMConfigType := reflect.TypeOf(*op)

	for i := 0; i < PMConfigValue.NumField(); i++ {
		PMFieldName := PMConfigType.Field(i).Name
		FieldValue := ConfigValue.FieldByName(PMFieldName)
		PMFieldValue := *(FieldValue.Interface().(*string))
		PMFieldValue = re.ReplaceAllString(PMFieldValue, "")   //去除非法字符
		PMFieldValue = rex.ReplaceAllString(PMFieldValue, "/") //去除多余斜杠
		PMConfigValue.Field(i).Set(reflect.ValueOf(&PMFieldValue))
	}
}

var re, _ = regexp.Compile("[^0-9a-zA-Z/]*")
var rex, _ = regexp.Compile("/+")
