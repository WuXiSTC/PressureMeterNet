package option

import "github.com/yindaheng98/gogisnet/grpc/server"

type Option struct {
	PressureMeterConfig PressureMeterConfig `yaml:"PressureMeterConfig" usage:"Option for PressureMeter."`
	GogisnetOption      server.Option       `yaml:"GogisnetOption" usage:"Option for gogisnet."`
}

func DefaultOption() Option {
	return Option{
		PressureMeterConfig: DefaultPressureMeterConfig(),
		GogisnetOption:      server.DefaultOption(),
	}
}
