package main

import (
	"PressureMeterMaster/option"
	"context"
	"gitee.com/WuXiSTC/PressureMeter"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"github.com/yindaheng98/gogisnet/grpc"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
	"github.com/yindaheng98/gogisnet/message"
	"net"
	"sync"
)

func ServerInit(opt option.Option) *server.Server {
	s := grpc.NewServer(opt.ServerInfoOption, opt.GogisnetOption) //Gogisnet初始化
	AddrSet := map[string]net.TCPAddr{}
	AddrSetMu := new(sync.Mutex)
	AddrList := func() []net.TCPAddr {
		addrList := make([]net.TCPAddr, len(AddrSet))
		i := 0
		for _, Addr := range AddrSet {
			addrList[i] = Addr
		}
		return addrList
	}
	s.Events.ClientNewConnection.AddHandler(func(info message.ClientInfo) {
		if Addr, err := net.ResolveTCPAddr("", string(info.(*pb.ClientInfo).AdditionalInfo)); err == nil {
			AddrSetMu.Lock()
			defer AddrSetMu.Unlock()
			AddrSet[Addr.String()] = *Addr
			*opt.PressureMeterConfig.ModelConfig.TaskConfig.IPList = AddrList() //修改地址表
		}
	})
	s.Events.ClientNewConnection.Enable()
	s.Events.ClientDisconnection.AddHandler(func(info message.ClientInfo) {
		if Addr, err := net.ResolveTCPAddr("", string(info.(*pb.ClientInfo).AdditionalInfo)); err == nil {
			AddrSetMu.Lock()
			defer AddrSetMu.Unlock()
			delete(AddrSet, Addr.String())
			*opt.PressureMeterConfig.ModelConfig.TaskConfig.IPList = AddrList() //修改地址表
		}
	})
	s.Events.ClientDisconnection.Enable()
	return s
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
