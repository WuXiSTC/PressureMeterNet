package option

import "github.com/yindaheng98/gogisnet/grpc/server"

type ListenerOption struct {
	S2SListenNetwork        *string
	S2SListenAddr           *string
	S2CListenNetwork        *string
	S2CListenAddr           *string
	GraphQueryListenNetwork *string
	GraphQueryListenAddr    *string
}

func DefaultListenerOption() ListenerOption {
	o := server.DefaultListenerOption()
	return ListenerOption{
		S2SListenNetwork:        &o.S2SListenNetwork,
		S2SListenAddr:           &o.S2SListenAddr,
		S2CListenNetwork:        &o.S2CListenNetwork,
		S2CListenAddr:           &o.S2CListenAddr,
		GraphQueryListenNetwork: &o.GraphQueryListenNetwork,
		GraphQueryListenAddr:    &o.GraphQueryListenAddr,
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
