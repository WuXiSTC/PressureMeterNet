package option

import (
	"github.com/yindaheng98/gogisnet/grpc/option"
	"github.com/yindaheng98/gogisnet/grpc/server"
)

type Option struct {
	ServerInfoOption    option.ServerInfoOption `yaml:"ServerInfoOption" usage:"Information about this server."`
	PressureMeterConfig PressureMeterConfig     `yaml:"PressureMeterConfig" usage:"Option for PressureMeter."`
	GogisnetOption      server.Option           `yaml:"GogisnetOption" usage:"Option for gogisnet."`
	ListenerOption      ListenerOption          `yaml:"ListenerOption" usage:"Option for port listen."`
	AccessAddr          string                  `yaml:"AccessAddr" usage:"The addr for the accessing from outside network."`
}

func DefaultOption() Option {
	return Option{
		ServerInfoOption:    option.DefaultServerInfoOption(),
		PressureMeterConfig: DefaultPressureMeterConfig(),
		GogisnetOption:      server.DefaultOption(),
		ListenerOption:      defaultListenerOption(),
		AccessAddr:          option.GetIP() + ":8080",
	}
}

type ListenerOption struct {
	GogisnetListenerOption  option.ListenerOption `yaml:"GogisnetListenerOption" usage:"Listener option for gogisnet."`
	PressureMeterListenAddr string                `yaml:"PressureMeterListenAddr" usage:"Listen address of the iris server in PressureMeter."`
}

func defaultListenerOption() ListenerOption {
	return ListenerOption{
		GogisnetListenerOption:  option.DefaultListenerOption(),
		PressureMeterListenAddr: "0.0.0.0:8080",
	}
}
