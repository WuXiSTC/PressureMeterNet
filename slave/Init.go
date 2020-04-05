package main

import (
	"PressureMeterNet/slave/option"
	"fmt"
	"github.com/yindaheng98/go-utility/MultiProcessController"
	"github.com/yindaheng98/gogisnet/grpc/client"
	"log"
)

func ClientInit(opt option.Option, mpc *MultiProcessController.MultiProcessController) *client.Client {
	c := client.New(opt.ClientInfoOption, opt.GogisnetOption)
	EventInit(c, mpc, opt.PressureMeterConfig)
	return c
}

type Logger struct {
}

func (l Logger) Log(args ...interface{}) {
	log.Println(append([]interface{}{"jmeter-server-->"}, args...)...)
}

type Commander struct {
	basePort uint16
}

func (c Commander) GenerateStartCommand(i uint64) (name string, args []string) {
	return "jmeter", []string{
		fmt.Sprintf("-Jserver.rmi.localport=%d", c.basePort+uint16(i)),
		fmt.Sprintf("-Dserver_port=%d", c.basePort+uint16(i)),
		"--server", "-Jserver.rmi.ssl.disable=true"}
}

func (c Commander) GenerateStopCommand(i uint64) (name string, args []string) {
	return "sh", []string{"-c",
		fmt.Sprintf("kill $(echo `ps -ef | grep %d | awk '{print $1}'`)", c.basePort+uint16(i)),
	}
}

func PressureMeterInit(opt option.PressureMeterConfig) *MultiProcessController.MultiProcessController {
	return MultiProcessController.New(Commander{opt.BaseListenPort}, Logger{})
}
