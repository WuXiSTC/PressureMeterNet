package option

import "github.com/yindaheng98/gogisnet/grpc/server"

type ListenerOption struct {
	S2SListenNetwork        *string `yaml:"S2SListenNetwork" usage:"Network type that the S2SRegistry listen on (e.g. tcp, udp)."`
	S2SListenAddr           *string `yaml:"S2SListenAddr" usage:"Network address and port that the S2SRegistry listen on."`
	S2CListenNetwork        *string `yaml:"S2CListenNetwork" usage:"Network type that the S2CRegistry listen on (e.g. tcp, udp)."`
	S2CListenAddr           *string `yaml:"S2CListenAddr" usage:"Network address and port that the S2CRegistry listen on."`
	GraphQueryListenNetwork *string `yaml:"GraphQueryListenNetwork" usage:"Network type that the GraphQueryRegistry listen on (e.g. tcp, udp)."`
	GraphQueryListenAddr    *string `yaml:"GraphQueryListenAddr" usage:"Network address and port that the GraphQuery server listen on."`
	PressureMeterListenAddr *string `yaml:"PressureMeterListenAddr" usage:"Network address and port that the PressureMeter server listen on."`
}

func DefaultListenerOption() ListenerOption {
	o := server.DefaultListenerOption()
	PressureMeterListenAddr := "localhost:8080"
	return ListenerOption{
		S2SListenNetwork:        &o.S2SListenNetwork,
		S2SListenAddr:           &o.S2SListenAddr,
		S2CListenNetwork:        &o.S2CListenNetwork,
		S2CListenAddr:           &o.S2CListenAddr,
		GraphQueryListenNetwork: &o.GraphQueryListenNetwork,
		GraphQueryListenAddr:    &o.GraphQueryListenAddr,
		PressureMeterListenAddr: &PressureMeterListenAddr,
	}
}
func (o ListenerOption) PutOption(op *server.ListenerOption) {
	op.S2SListenNetwork = *o.S2SListenNetwork
	op.S2SListenAddr = *o.S2SListenAddr
	op.S2CListenNetwork = *o.S2CListenNetwork
	op.S2CListenAddr = *o.S2CListenAddr
	op.GraphQueryListenNetwork = *o.GraphQueryListenNetwork
	op.GraphQueryListenAddr = *o.GraphQueryListenAddr
}
