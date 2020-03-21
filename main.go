package main

import (
	"PressureMeterMaster/option"
	"context"
	"gitee.com/WuXiSTC/PressureMeter"
	"github.com/kataras/iris"
	"github.com/yindaheng98/gogisnet/grpc"
	"log"
	"time"
)

func main() {
	opt, exit := option.Generate(func(i ...interface{}) {
		log.Println(i...)
	})
	if exit { //如果要退出
		return //就直接退出
	}
	//fmt.Println(opt)
	server := grpc.NewServer(opt.ServerInfoOption, opt.GogisnetOption) //Gogisnet初始化

	ctxBackground := context.Background()
	ctx, cancel := context.WithCancel(ctxBackground)
	defer cancel()
	PressureMeterConfig := PressureMeter.Config{}
	opt.PressureMeterConfig.PutConfig(&PressureMeterConfig)
	app := PressureMeter.Init(ctx, PressureMeterConfig) //PressureMeter服务器初始化

	//TODO：获取Gogisnet全网连接图API

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancelInterrupt := context.WithTimeout(ctxBackground, timeout)
		defer cancelInterrupt()
		_ = app.Shutdown(ctx) //关闭所有主机
		cancel()
	})

	errChan := make(chan error)
	go func() {
		errChan <- server.Run(ctx, opt.ListenerOption.GogisnetListenerOption)
	}()
	go func() {
		errChan <- app.Run(iris.Addr(opt.ListenerOption.PressureMeterListenAddr), iris.WithoutServerError(iris.ErrServerClosed))
	}()
	if err := <-errChan; err != nil {
		log.Println(err)
		cancel()
	}
}
