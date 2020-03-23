package main

import (
	"PressureMeterMaster/hsync"
	"PressureMeterMaster/option"
	"context"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"github.com/yindaheng98/gogisnet/grpc"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
	"github.com/yindaheng98/gogisnet/message"
)

func ServerInit(opt option.Option) (*server.Server, *hsync.ClientIPSync) {
	s := grpc.NewServer(opt.ServerInfoOption, opt.GogisnetOption) //Gogisnet初始化
	sync := hsync.NewClientIPSync(opt.HostsFileSyncOption.HostsFilePath,
		func(u uint64) string {
			return fmt.Sprintf(opt.HostsFileSyncOption.HostsFormat, u)
		})
	s.Events.ClientNewConnection.AddHandler(func(info message.ClientInfo) {
		sync.AddClient(info.GetClientID(), string(info.(*pb.ClientInfo).AdditionalInfo))
	})
	s.Events.ClientNewConnection.Enable()
	s.Events.ClientDisconnection.AddHandler(func(info message.ClientInfo) {
		sync.DelClient(info.GetClientID())
	})
	s.Events.ClientDisconnection.Enable()
	return s, sync
}

func PressureMeterInit(ctx context.Context, opt option.PressureMeterConfig) *iris.Application {
	PressureMeterConfig := PressureMeter.Config{}
	opt.PutConfig(&PressureMeterConfig)
	return PressureMeter.Init(ctx, PressureMeterConfig) //PressureMeter服务器初始化
}

func GraphAPIInit(s *server.Server, app *iris.Application, url string) {
	app.Get(url, func(ctx irisContext.Context) {
		_, _ = ctx.Write([]byte(s.GetGraph(context.Background()).String()))
	})
}
