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
	HostsFileSyncOption HostsFileSyncOption     `yaml:"HostsFileSyncOption" usage:"Option for hosts file sync."`
}

func DefaultOption() Option {
	return Option{
		ServerInfoOption:    option.DefaultServerInfoOption(),
		PressureMeterConfig: DefaultPressureMeterConfig(),
		GogisnetOption:      server.DefaultOption(),
		ListenerOption:      defaultListenerOption(),
		HostsFileSyncOption: defaultHostsFileSyncOption(),
	}
}

type ListenerOption struct {
	GogisnetListenerOption  option.ListenerOption `yaml:"GogisnetListenerOption" usage:"Listener option for gogisnet."`
	PressureMeterListenAddr string                `yaml:"PressureMeterListenAddr" usage:"Listen address of the iris server in PressureMeter."`
}

func defaultListenerOption() ListenerOption {
	return ListenerOption{
		GogisnetListenerOption:  option.DefaultListenerOption(),
		PressureMeterListenAddr: ":80",
	}
}

type HostsFileSyncOption struct {
	HostsFilePath string `yaml:"HostsFilePath" usage:"Path to your hosts file."`
	HostsFormat   string `yaml:"HostsFormat" usage:"Format of the host names."`
}

func defaultHostsFileSyncOption() HostsFileSyncOption {
	return HostsFileSyncOption{
		HostsFilePath: "/etc/hosts",
		HostsFormat:   "jmeter%02d.test.wxstc.org.cn",
	}
}
