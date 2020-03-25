package main

import (
	"PressureMeterNet/master/option"
	"context"
	"github.com/kataras/iris"
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
	server := ServerInit(opt) //Gogisnet初始化

	ctxBackground := context.Background()
	ctx, cancel := context.WithCancel(ctxBackground)
	defer cancel()

	app := PressureMeterInit(ctx, opt.PressureMeterConfig) //PressureMeter服务器初始化

	GraphAPIInit(server, app, opt.GraphQueryURL) //全网连接图API初始化

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
