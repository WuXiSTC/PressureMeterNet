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

type Option struct {
	InitServerOption InitServerOption
	ServiceOption    ServiceOption
	GRPCOption       GRPCOption
	BoardCastOption  BoardCastOption
}

func DefaultOption() Option {
	return Option{
		InitServerOption: DefaultInitServerOption(),
		ServiceOption:    DefaultServiceOption(),
		GRPCOption:       DefaultGRPCOption(),
		BoardCastOption:  DefaultBoardCastOption(),
	}
}
func (o Option) PutOption() (op server.Option) {
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
	S2SBoardCastAddr        *string
	S2CBoardCastAddr        *string
	GraphQueryBoardCastAddr *string
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
	ServerID    *string
	ServiceType *string
	S2SAddress  *string
	S2CAddress  *string
}

func DefaultInitServerOption() InitServerOption {
	ServerID := "default"
	ServiceType := "default"
	S2SAddress := "default"
	S2CAddress := "default"
	return InitServerOption{
		ServerID:    &ServerID,
		ServiceType: &ServiceType,
		S2SAddress:  &S2SAddress,
		S2CAddress:  &S2CAddress,
	}
}
func (o InitServerOption) PutOption() *pb.S2SInfo {
	ServerInfo := &pb.ServerInfo{
		ServerID:    *o.ServerID,
		ServiceType: *o.ServiceType,
	}
	return &pb.S2SInfo{
		ServerInfo:         ServerInfo,
		ResponseSendOption: &pb.ResponseSendOption{},
		RequestSendOption:  &pb.RequestSendOption{Addr: *o.S2SAddress},
		S2CInfo: &pb.S2CInfo{
			ServerInfo:        ServerInfo,
			RequestSendOption: &pb.RequestSendOption{Addr: *o.S2CAddress},
		},
	}
}

type ServiceOption struct {
	S2SRegistryOption   RegistryOption
	S2SRegistrantOption RegistrantOption
	S2CRegistryOption   RegistryOption
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
	MaxRegistrants             *uint64
	LogTimeoutControllerOption LogTimeoutControllerOption
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
	op.TimeoutController = o.LogTimeoutControllerOption.PutOption()
}

type LogTimeoutControllerOption struct {
	MinimumTime    *uint64
	MaximumTime    *uint64
	IncreaseFactor *float64
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
func (o LogTimeoutControllerOption) PutOption() *TimeoutController.LogTimeoutController {
	return TimeoutController.NewLogTimeoutController(
		time.Duration(*o.MinimumTime), time.Duration(*o.MaximumTime),
		*o.IncreaseFactor)
}

type RegistrantOption struct {
	RegistryN *uint64
}

func DefaultRegistrantOption() RegistrantOption {
	RegistryN := uint64(4)
	return RegistrantOption{RegistryN: &RegistryN}
}
func (o RegistrantOption) PutOption(op *service.RegistrantOption) {
	op.RegistryN = *o.RegistryN
}

type GRPCOption struct {
	S2SRegistryOption   GRPCRegistryOption
	S2CRegistryOption   GRPCRegistryOption
	S2SRegistrantOption GRPCRegistrantOption
	CandidateListOption CandidateListOption
	GraphQueryOption    GraphQueryOption
}

func DefaultGRPCOption() GRPCOption {
	return GRPCOption{
		S2SRegistryOption:   DefaultGRPCRegistryOption(),
		S2CRegistryOption:   DefaultGRPCRegistryOption(),
		S2SRegistrantOption: DefaultGRPCRegistrantOption(),
		CandidateListOption: DefaultCandidateListOption(),
		GraphQueryOption:    DefaultGraphQueryOption(),
	}
}
func (o GRPCOption) PutOption(op *server.GRPCOption) {
	o.S2SRegistryOption.PutOption(&op.S2SRegistryOption)
	o.S2CRegistryOption.PutOption(&op.S2CRegistryOption)
	o.S2SRegistrantOption.PutOption(&op.S2SRegistrantOption)
	o.CandidateListOption.PutOption(&op.CandidateListOption)
	o.GraphQueryOption.PutOption(&op.GraphQueryOption)
}

type GRPCRegistryOption struct {
	BufferSize *uint64
}

func DefaultGRPCRegistryOption() GRPCRegistryOption {
	BufferSize := uint64(100)
	return GRPCRegistryOption{BufferSize: &BufferSize}
}
func (o GRPCRegistryOption) PutOption(op *registry.GRPCRegistryOption) {
	op.BufferLen = *o.BufferSize
}

type GRPCRegistrantOption struct {
	MaxDialHoldDuration *uint64
}

func DefaultGRPCRegistrantOption() GRPCRegistrantOption {
	MaxDialHoldDuration := uint64(100e9)
	return GRPCRegistrantOption{MaxDialHoldDuration: &MaxDialHoldDuration}
}
func (o GRPCRegistrantOption) PutOption(op *registrant.GRPCRegistrantOption) {
	op.MaxDialHoldDuration = time.Duration(*o.MaxDialHoldDuration)
}

type CandidateListOption struct {
	InitTimeLimit  *uint64
	InitRetryLimit *uint64
	Size           *uint64
	PingTimeLimit  *uint64
}

func DefaultCandidateListOption() CandidateListOption {
	InitTimeLimit := uint64(1e9)
	InitRetryLimit := uint64(10)
	Size := uint64(16)
	PingTimeLimit := uint64(2e9)
	return CandidateListOption{
		InitTimeLimit:  &InitTimeLimit,
		InitRetryLimit: &InitRetryLimit,
		Size:           &Size,
		PingTimeLimit:  &PingTimeLimit,
	}
}
func (o CandidateListOption) PutOption(op *registrant.CandidateListOption) {
	op.InitTimeout = time.Duration(*o.InitTimeLimit)
	op.InitRetryN = *o.InitRetryLimit
	op.Size = *o.Size
	op.MaxPingTimeout = time.Duration(*o.PingTimeLimit)
}

type GraphQueryOption struct {
	GraphQueryClientOption GRPCRegistrantOption
}

func DefaultGraphQueryOption() GraphQueryOption {
	return GraphQueryOption{GraphQueryClientOption: DefaultGRPCRegistrantOption()}
}
func (o GraphQueryOption) PutOption(op *graph.GraphQueryOption) {
	o.GraphQueryClientOption.PutOption((*registrant.GRPCRegistrantOption)(&op.GraphQueryClientOption))
}
