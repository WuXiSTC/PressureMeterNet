package main

import (
	"PressureMeterMaster/option"
	"fmt"
)

func main() {
	opt, exit := option.Generate(func(i ...interface{}) {
		fmt.Println(i...)
	})
	if exit { //如果要退出
		return //就直接退出
	}
	fmt.Println(opt)
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
