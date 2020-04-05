package option

import (
	"github.com/yindaheng98/gogisnet/grpc/client"
	"github.com/yindaheng98/gogisnet/grpc/option"
)

type Option struct {
	ClientInfoOption    option.ClientInfoOption
	GogisnetOption      client.Option
	PressureMeterConfig PressureMeterConfig
}

func DefaultOption() Option {
	return Option{
		ClientInfoOption:    option.DefaultClientInfoOption(),
		GogisnetOption:      client.DefaultOption(),
		PressureMeterConfig: defaultPressureMeterConfig(),
	}
}

type PressureMeterConfig struct {
	BoardCastAddr  string
	BaseListenPort uint16
}

func defaultPressureMeterConfig() PressureMeterConfig {
	return PressureMeterConfig{
		BoardCastAddr:  option.GetIP(),
		BaseListenPort: 1099,
	}
}
