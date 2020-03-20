package main

import (
	"PressureMeterMaster/option"
	"fmt"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
)

func main() {
	opt := option.GenerateOption(func(i ...interface{}) {
		fmt.Println(i...)
	})
	ServerInfoOption := opt.ServerInfo
	GogisnetOption := opt.GogisnetOption
	ListenerOption := opt.GRPCListenerOption
	PressureMeterConfig := opt.PressureMeterConfig
	ServerInfo := new(pb.ServerInfo)
	ServerInfoOption.PutOption(ServerInfo)
	s := grpc.NewServer(ServerInfo, GogisnetOption.PutOption())
	/*
		ServerInfoOption, GogisnetOption, ListenerOption := option.GenerateOption()
		ServerInfo := new(pb.ServerInfo)
		ServerInfoOption.PutOption(ServerInfo)
		s := grpc.NewServer(ServerInfo, GogisnetOption.PutOption())

		ctx := context.Background()
		a:=PressureMeter.Init(ctx,PressureMeter.Config{
			ModelConfig:  Model.Config{},
			URLConfig:    PressureMeter.URLConfig{},
			LoggerConfig: nil,
		})
		Listener := new(server.ListenerOption)
		ListenerOption.PutOption(Listener)
		err:=s.Run(ctx, *Listener)
	*/
}
