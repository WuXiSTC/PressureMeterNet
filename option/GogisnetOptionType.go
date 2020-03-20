package option

import (
	"github.com/yindaheng98/gogisnet/grpc/protocol/graph"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/protocol/registrant"
	"github.com/yindaheng98/gogisnet/grpc/protocol/registry"
	"github.com/yindaheng98/gogisnet/grpc/server"
	service "github.com/yindaheng98/gogisnet/server"
	"github.com/yindaheng98/gogistry/example/TimeoutController"
	"time"
)

type GogisnetOption struct {
	InitServerOption InitServerOption `yaml:"InitServerOption" usage:"Information for the first server that will be connected by current server."`
	ServiceOption    ServiceOption    `yaml:"ServiceOption" usage:"Option for the service in gogisnet."`
	GRPCOption       GRPCOption       `yaml:"GRPCOption" usage:"Option for the gRPC server in gogisnet."`
	BoardCastOption  BoardCastOption  `yaml:"BoardCastOption" usage:"External access addresses for other servers (board cast with messages) to find the current server."`
}

func DefaultGogisnetOption() GogisnetOption {
	return GogisnetOption{
		InitServerOption: DefaultInitServerOption(),
		ServiceOption:    DefaultServiceOption(),
		GRPCOption:       DefaultGRPCOption(),
		BoardCastOption:  DefaultBoardCastOption(),
	}
}
func (o GogisnetOption) PutOption(op *server.Option) {
	op, err := server.DefaultOption(
		*o.BoardCastOption.S2SBoardCastAddr,
		*o.BoardCastOption.S2CBoardCastAddr,
		*o.BoardCastOption.GraphQueryBoardCastAddr,
		o.InitServerOption.PutOption())
	if err != nil {
		panic(err)
	}
	return
}

type BoardCastOption struct {
	S2SBoardCastAddr        *string `yaml:"S2SBoardCastAddr" usage:"External access address for other gogisnet servers."`
	S2CBoardCastAddr        *string `yaml:"S2CBoardCastAddr" usage:"External access address for gogisnet clients."`
	GraphQueryBoardCastAddr *string `yaml:"GraphQueryBoardCastAddr" usage:"External access address of GraphQuery service."`
}

func DefaultBoardCastOption() BoardCastOption {
	IPAddr := GetIP()
	S2SBoardCastAddr := IPAddr + ":4241"
	S2CBoardCastAddr := IPAddr + ":4240"
	GraphQueryBoardCastAddr := IPAddr + ":4242"
	return BoardCastOption{
		S2SBoardCastAddr:        &S2SBoardCastAddr,
		S2CBoardCastAddr:        &S2CBoardCastAddr,
		GraphQueryBoardCastAddr: &GraphQueryBoardCastAddr,
	}
}

type InitServerOption struct {
	ServerID    *string `yaml:"ServerID" usage:"Unique ID for the identification of the server."`
	ServiceType *string `yaml:"ServiceType" usage:"Type of the server (when connecting, the type of client and the server must be the same)."`
	S2SAddress  *string `yaml:"S2SAddress" usage:"External access address for other gogisnet servers."`
	S2CAddress  *string `yaml:"S2CAddress" usage:"External access address for gogisnet clients."`
}

func DefaultInitServerOption() InitServerOption {
	ServerID := "undefined"
	ServiceType := "undefined"
	S2SAddress := "undefined"
	S2CAddress := "undefined"
	return InitServerOption{
		ServerID:    &ServerID,
		ServiceType: &ServiceType,
		S2SAddress:  &S2SAddress,
		S2CAddress:  &S2CAddress,
	}
}
func (o InitServerOption) PutOption(op *pb.S2SInfo) {
	op.ServerInfo.ServerID = *o.ServerID
	op.ServerInfo.ServiceType = *o.ServiceType
	op.RequestSendOption.Addr = *o.S2SAddress
	op.S2CInfo.ServerInfo.ServerID = *o.ServerID
	op.S2CInfo.ServerInfo.ServiceType = *o.ServiceType
	op.S2CInfo.RequestSendOption.Addr = *o.S2CAddress
}

type ServiceOption struct {
	S2SRegistryOption   RegistryOption   `yaml:"S2SRegistryOption" usage:"Option for S2SRegistry in gogisnet server."`
	S2SRegistrantOption RegistrantOption `yaml:"S2SRegistrantOption" usage:"Option for S2SRegistrant in gogisnet server."`
	S2CRegistryOption   RegistryOption   `yaml:"S2CRegistryOption" usage:"Option for S2CRegistry in gogisnet server."`
}

func DefaultServiceOption() ServiceOption {
	return ServiceOption{
		S2SRegistryOption:   DefaultRegistryOption(),
		S2SRegistrantOption: DefaultRegistrantOption(),
		S2CRegistryOption:   DefaultRegistryOption(),
	}
}
func (o ServiceOption) PutOption(op *service.Option) {
	o.S2SRegistryOption.PutOption(&op.S2SRegistryOption)
	o.S2SRegistrantOption.PutOption(&op.S2SRegistrantOption)
	o.S2CRegistryOption.PutOption(&op.S2CRegistryOption)
}

type RegistryOption struct {
	MaxRegistrants             *uint64                    `yaml:"MaxRegistrants" usage:"The number of registrants that a registry can connect to at most."`
	LogTimeoutControllerOption LogTimeoutControllerOption `yaml:"LogTimeoutControllerOption" usage:"Option for the TimeoutController in registry."`
}

func DefaultRegistryOption() RegistryOption {
	MaxRegistrants := uint64(4)
	return RegistryOption{
		MaxRegistrants:             &MaxRegistrants,
		LogTimeoutControllerOption: DefaultLogTimeoutControllerOption(),
	}
}
func (o RegistryOption) PutOption(op *service.RegistryOption) {
	op.MaxRegistrants = *o.MaxRegistrants
	o.LogTimeoutControllerOption.PutOption(op.TimeoutController.(*TimeoutController.LogTimeoutController))
}

type LogTimeoutControllerOption struct {
	MinimumTime    *uint64  `yaml:"MinimumTime" usage:"Minimum output of the TimeoutController (timeout(0)=MinimumTime)."`
	MaximumTime    *uint64  `yaml:"MaximumTime" usage:"Maximum output of the TimeoutController."`
	IncreaseFactor *float64 `yaml:"IncreaseFactor" usage:"IncreaseFactor of the TimeoutController (timeout(n)=timeout(n-1)+(MaximumTime-timeout(n-1))/IncreaseFactor)."`
}

func DefaultLogTimeoutControllerOption() LogTimeoutControllerOption {
	MinimumTime := uint64(1e9)
	MaximumTime := uint64(10e9)
	IncreaseFactor := float64(2)
	return LogTimeoutControllerOption{
		MinimumTime:    &MinimumTime,
		MaximumTime:    &MaximumTime,
		IncreaseFactor: &IncreaseFactor,
	}
}
func (o LogTimeoutControllerOption) PutOption(op *TimeoutController.LogTimeoutController) {
	op.MinimumTime = time.Duration(*o.MinimumTime)
	op.MaximumTime = time.Duration(*o.MaximumTime)
	op.IncreaseFactor = *o.IncreaseFactor
}

type RegistrantOption struct {
	RegistryN *uint64 `yaml:"RegistryN" usage:"The number of registries that a registrant can connect at most."`
}

func DefaultRegistrantOption() RegistrantOption {
	RegistryN := uint64(4)
	return RegistrantOption{RegistryN: &RegistryN}
}
func (o RegistrantOption) PutOption(op *service.RegistrantOption) {
	op.RegistryN = *o.RegistryN
}

type GRPCOption struct {
	S2SRegistryOption   GRPCRegistryOption   `yaml:"S2SRegistryOption" usage:"Option for gRPC server in S2SRegistry."`
	S2CRegistryOption   GRPCRegistryOption   `yaml:"S2CRegistryOption" usage:"Option for gRPC server in S2CRegistry."`
	S2SRegistrantOption GRPCRegistrantOption `yaml:"S2SRegistrantOption" usage:"Option for gRPC client in S2SRegistrant."`
	GraphQueryOption    GraphQueryOption     `yaml:"GraphQueryOption" usage:"Option for gRPC server in GraphQuery service."`
}

func DefaultGRPCOption() GRPCOption {
	return GRPCOption{
		S2SRegistryOption:   DefaultGRPCRegistryOption(),
		S2CRegistryOption:   DefaultGRPCRegistryOption(),
		S2SRegistrantOption: DefaultGRPCRegistrantOption(),
		GraphQueryOption:    DefaultGraphQueryOption(),
	}
}
func (o GRPCOption) PutOption(op *server.GRPCOption) {
	o.S2SRegistryOption.PutOption(&op.S2SRegistryOption)
	o.S2CRegistryOption.PutOption(&op.S2CRegistryOption)
	o.S2SRegistrantOption.PutOption(&op.S2SRegistrantOption)
	o.GraphQueryOption.PutOption(&op.GraphQueryOption)
}

type GRPCRegistryOption struct {
	BufferSize *uint64 `yaml:"BufferSize" usage:"Size of the request buffer in gRPC server."`
}

func DefaultGRPCRegistryOption() GRPCRegistryOption {
	BufferSize := uint64(100)
	return GRPCRegistryOption{BufferSize: &BufferSize}
}
func (o GRPCRegistryOption) PutOption(op *registry.GRPCRegistryOption) {
	op.BufferLen = *o.BufferSize
}

type GRPCRegistrantOption struct {
	MaxDialHoldDuration *uint64 `yaml:"MaxDialHoldDuration" usage:"Time that a gRPC dial hold in gRPC client at most."`
}

func DefaultGRPCRegistrantOption() GRPCRegistrantOption {
	MaxDialHoldDuration := uint64(100e9)
	return GRPCRegistrantOption{MaxDialHoldDuration: &MaxDialHoldDuration}
}
func (o GRPCRegistrantOption) PutOption(op *registrant.GRPCRegistrantOption) {
	op.MaxDialHoldDuration = time.Duration(*o.MaxDialHoldDuration)
}

type CandidateListOption struct {
	DefaultTimeLimit  *uint64 `yaml:"InitTimeLimit" usage:"The first timeout value returned by CandidateList."`
	DefaultRetryLimit *uint64 `yaml:"InitRetryLimit" usage:"The first retryN value returned by CandidateList."`
	Size              *uint64 `yaml:"Size" usage:"The number of server that the CandidateList can store at most."`
	PingTimeLimit     *uint64 `yaml:"PingTimeLimit" usage:"Time limit of the 'PING' operation."`
}

func DefaultCandidateListOption() CandidateListOption {
	DefaultTimeLimit := uint64(1e9)
	DefaultRetryLimit := uint64(10)
	Size := uint64(16)
	PingTimeLimit := uint64(2e9)
	return CandidateListOption{
		DefaultTimeLimit:  &DefaultTimeLimit,
		DefaultRetryLimit: &DefaultRetryLimit,
		Size:              &Size,
		PingTimeLimit:     &PingTimeLimit,
	}
}
func (o CandidateListOption) PutOption(op *registrant.PingerCandidateListOption) {
	op.DefaultTimeout = time.Duration(*o.DefaultTimeLimit)
	op.DefaultTimeout = time.Duration(*o.DefaultRetryLimit)
	op.Size = *o.Size
	op.MaxPingTimeout = time.Duration(*o.PingTimeLimit)
}

type GraphQueryOption struct {
	GraphQueryClientOption GRPCRegistrantOption `yaml:"GraphQueryClientOption" usage:"Option for the gRPC client in GraphQuery service."`
}

func DefaultGraphQueryOption() GraphQueryOption {
	return GraphQueryOption{GraphQueryClientOption: DefaultGRPCRegistrantOption()}
}
func (o GraphQueryOption) PutOption(op *graph.GraphQueryOption) {
	o.GraphQueryClientOption.PutOption((*registrant.GRPCRegistrantOption)(&op.GraphQueryClientOption))
}
