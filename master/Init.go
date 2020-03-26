package main

import (
	"PressureMeterNet/master/option"
	"context"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"github.com/yindaheng98/gogisnet/grpc"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
	"github.com/yindaheng98/gogisnet/message"
	"log"
	"net"
	"sync"
)

func ResolveAddr(addr string) (*net.TCPAddr, error) {
	if Addr, err := net.ResolveTCPAddr("", addr); err == nil {
		return Addr, nil
	}
	if IP, err := net.ResolveIPAddr("", addr); err == nil {
		return &net.TCPAddr{
			IP:   IP.IP,
			Port: 1099,
			Zone: IP.Zone,
		}, nil
	} else {
		return nil, err
	}
}

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
	UpdateClient := func(info message.ClientInfo) {
		AddrSetMu.Lock()
		defer AddrSetMu.Unlock()
		addr := string(info.(*pb.ClientInfo).AdditionalInfo)
		if Addr, err := ResolveAddr(addr); err == nil {
			AddrSet[info.GetClientID()] = *Addr
			log.Println(fmt.Sprintf("Client '%s' updated: %s", info.GetClientID(), Addr.String()))
		} else {
			delete(AddrSet, info.GetClientID())
			log.Println(fmt.Sprintf("Client '%s' update failed: %s\n", info.GetClientID(), addr), err)
		}
		*opt.PressureMeterConfig.ModelConfig.TaskConfig.IPList = AddrList() //修改地址表
	}
	DeleteClient := func(info message.ClientInfo) {
		AddrSetMu.Lock()
		defer AddrSetMu.Unlock()
		delete(AddrSet, info.GetClientID())
		log.Println(fmt.Sprintf("Client '%s' deleted", info.GetClientID()))
		*opt.PressureMeterConfig.ModelConfig.TaskConfig.IPList = AddrList() //修改地址表
	}
	s.Events.ClientNewConnection.AddHandler(UpdateClient)
	s.Events.ClientNewConnection.Enable()
	s.Events.S2CRegistryEvent.UpdateConnection.AddHandler(func(info message.C2SInfo) { UpdateClient(info.ClientInfo) })
	s.Events.ClientDisconnection.AddHandler(DeleteClient)
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
