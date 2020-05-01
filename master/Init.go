package main

import (
	"PressureMeterNet/master/option"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter"
	"gitee.com/WuXiSTC/PressureMeter/Model/TaskList"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"github.com/yindaheng98/gogisnet/grpc"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
)

func ServerInit(opt option.Option) *server.Server {
	opt.ServerInfoOption.AdditionalInfo["AccessAddress"] = opt.AccessAddr
	opt.ServerInfoOption.AdditionalInfo["TaskAccN"] =
		fmt.Sprintf("%d", opt.PressureMeterConfig.ModelConfig.DaemonConfig.TaskAccN)
	s := grpc.NewServer(opt.ServerInfoOption, opt.GogisnetOption) //Gogisnet初始化
	EventInit(s, opt.PressureMeterConfig.ModelConfig.DaemonConfig.TaskAccN)
	return s
}

func PressureMeterInit(s *server.Server, ctx context.Context, opt option.PressureMeterConfig) *iris.Application {
	PressureMeterConfig := PressureMeter.Config{}
	opt.PutConfig(&PressureMeterConfig)
	PressureMeterConfig.ModelConfig.UpdateStateCallback = func(list TaskList.TaskStateList) {
		si := s.GetS2SInfo().ServerInfo.(*pb.ServerInfo)
		tl, _ := json.Marshal(list)
		si.AdditionalInfo["TaskList"] = string(tl)
	}
	return PressureMeter.Init(ctx, PressureMeterConfig) //PressureMeter服务器初始化
}

func GraphAPIInit(s *server.Server, app *iris.Application, opt option.Option) {
	app.Get(opt.PressureMeterConfig.URLConfig.GraphQuery, func(ctx irisContext.Context) {
		_, _ = ctx.Write([]byte(s.GetGraph(context.Background()).String()))
	})
}
