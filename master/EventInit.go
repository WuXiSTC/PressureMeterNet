package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter/Model/Daemon"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
	"github.com/yindaheng98/gogisnet/message"
	"log"
	"net"
	"sync"
)

//IP地址批量转换
func ResolveAddrs(addrs []string) (Addrs []net.TCPAddr, err error) {
	Addrs = make([]net.TCPAddr, len(addrs))
	for i, addr := range addrs {
		if Addr, err := ResolveAddr(addr); err != nil {
			return nil, err
		} else {
			Addrs[i] = *Addr
		}
	}
	return
}

//IP地址转换
func ResolveAddr(addr string) (*net.TCPAddr, error) {
	if Addr, err := net.ResolveTCPAddr("", addr); err == nil {
		return Addr, nil
	}
	if IP, err := net.ResolveIPAddr("", addr); err == nil {
		return &net.TCPAddr{
			IP:   IP.IP,
			Port: 1099,
			Zone: IP.Zone,
		}, nil
	} else {
		return nil, err
	}
}

func EventInit(s *server.Server, TaskAccN uint16) {
	AddrSet := map[string][]net.TCPAddr{}
	AddrSetMu := new(sync.RWMutex)
	GetAddrList := func() [][]net.TCPAddr { //线程同步由调用此函数的操作保证
		addrList := make([][]net.TCPAddr, TaskAccN)
		for _, Addrs := range AddrSet { //遍历全部从机的地址列表
			if len(Addrs) >= int(TaskAccN) { //如果此从机的地址列表足够给每个并行任务分一份
				for TaskAccI := uint16(0); TaskAccI < TaskAccN; TaskAccI++ { //那就给每个并行任务分一份
					addrList[TaskAccI] = append(addrList[TaskAccI], Addrs[TaskAccI])
				}
			}
		}
		return addrList
	}
	handleErr := func(err error, info message.ClientInfo) {
		AdditionalInfo := make(map[string]string, len(info.(*pb.ClientInfo).AdditionalInfo))
		for k, v := range info.(*pb.ClientInfo).AdditionalInfo {
			AdditionalInfo[k] = string(v)
		}
		delete(AddrSet, info.GetClientID())
		log.Println(fmt.Sprintf("Client '%s' update failed: %s", info.GetClientID(), AdditionalInfo), err)
	}
	UpdateClient := func(info message.ClientInfo) {
		var addrs []string
		if err := json.Unmarshal(info.(*pb.ClientInfo).AdditionalInfo["JmeterAddresses"], &addrs); err != nil {
			handleErr(err, info) //反序列化出错
		} else if Addrs, err := ResolveAddrs(addrs); err != nil {
			handleErr(err, info) //IP出错
		} else if len(Addrs) < int(TaskAccN) { //长度不够
			handleErr(errors.New("len(Addrs) < int(TaskAccN),长度不足"), info)
		} else if err := Daemon.SetIPList(GetAddrList()); err != nil {
			handleErr(err, info) //修改地址表出错
		} else { //成功
			AddrSetMu.Lock()
			defer AddrSetMu.Unlock()
			AddrSet[info.GetClientID()] = Addrs
			log.Println(fmt.Sprintf("Client '%s' updated: ", info.GetClientID()), Addrs)
		}
	}
	DeleteClient := func(info message.ClientInfo) {
		AddrSetMu.Lock()
		defer AddrSetMu.Unlock()
		delete(AddrSet, info.GetClientID())
		log.Println(fmt.Sprintf("Client '%s' deleted", info.GetClientID()))
		if err := Daemon.SetIPList(GetAddrList()); err != nil {
			handleErr(err, info) //修改地址表出错
		}
	}
	s.Events.ClientNewConnection.AddHandler(UpdateClient)
	s.Events.ClientNewConnection.Enable()
	s.Events.S2CRegistryEvent.UpdateConnection.AddHandler(func(info message.C2SInfo) { UpdateClient(info.ClientInfo) })
	s.Events.S2CRegistryEvent.UpdateConnection.Enable()
	s.Events.ClientDisconnection.AddHandler(DeleteClient)
	s.Events.ClientDisconnection.Enable()
}
