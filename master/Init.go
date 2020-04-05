package main

import (
	"PressureMeterNet/master/option"
	"context"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"github.com/yindaheng98/gogisnet/grpc"
	"github.com/yindaheng98/gogisnet/grpc/server"
)

func ServerInit(opt option.Option) *server.Server {
	opt.ServerInfoOption.AdditionalInfo["AccessAddress"] = []byte(opt.AccessAddr)
	opt.ServerInfoOption.AdditionalInfo["TaskAccN"] =
		[]byte(fmt.Sprintf("%d", opt.PressureMeterConfig.ModelConfig.DaemonConfig.TaskAccN))
	s := grpc.NewServer(opt.ServerInfoOption, opt.GogisnetOption) //Gogisnet初始化
	EventInit(s, opt.PressureMeterConfig.ModelConfig.DaemonConfig.TaskAccN)
	return s
}

func PressureMeterInit(ctx context.Context, opt option.PressureMeterConfig) *iris.Application {
	PressureMeterConfig := PressureMeter.Config{}
	opt.PutConfig(&PressureMeterConfig)
	return PressureMeter.Init(ctx, PressureMeterConfig) //PressureMeter服务器初始化
}

func GraphAPIInit(s *server.Server, app *iris.Application, opt option.Option) {
	app.Get(opt.PressureMeterConfig.URLConfig.GraphQuery, func(ctx irisContext.Context) {
		_, _ = ctx.Write([]byte(s.GetGraph(context.Background()).String()))
	})
}
