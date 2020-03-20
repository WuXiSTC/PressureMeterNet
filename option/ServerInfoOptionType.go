package option

import (
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"math/rand"
	"time"
)

type ServerInfoOption struct {
	ServerID       *string `yaml:"ServerID" usage:"Unique ID for the identification of the server."`
	ServiceType    *string `yaml:"ServiceType" usage:"Type of the server (when connecting, the type of client and the server must be the same)."`
	AdditionalInfo *string `yaml:"AdditionalInfo" usage:"Some additional information you want to add to this server."`
}

func DefaultServerInfoOption() ServerInfoOption {
	ServerID := "SERVER-" + RandomString(64)
	ServiceType := "default"
	AdditionalInfo := ""
	return ServerInfoOption{
		ServerID:       &ServerID,
		ServiceType:    &ServiceType,
		AdditionalInfo: &AdditionalInfo,
	}
}
func (o ServerInfoOption) PutOption(op *pb.ServerInfo) {
	op.ServerID = *o.ServerID
	op.ServiceType = *o.ServiceType
}

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(src)
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
