package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/WuXiSTC/PressureMeter/Model/Task"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/grpc/server"
	"github.com/yindaheng98/gogisnet/message"
	"log"
	"sync"
)

func EventInit(s *server.Server, TaskAccN uint16) {
	AddrSet := map[string][]string{}
	AddrSetMu := new(sync.RWMutex)
	GetAddrList := func() [][]string { //线程同步由调用此函数的操作保证
		addrList := make([][]string, TaskAccN)
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
		var Addrs []string
		if err := json.Unmarshal([]byte(info.(*pb.ClientInfo).AdditionalInfo["JmeterAddresses"]), &Addrs); err != nil {
			handleErr(err, info) //反序列化出错
		} else if len(Addrs) < int(TaskAccN) { //长度不够
			handleErr(errors.New("len(Addrs) < int(TaskAccN),长度不足"), info)
		} else if err := Task.SetHosts(GetAddrList()); err != nil {
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
		if err := Task.SetHosts(GetAddrList()); err != nil {
			handleErr(err, info) //修改地址表出错
		}
	}
	s.Events.ClientNewConnection.AddHandler(UpdateClient)
	s.Events.ClientNewConnection.Enable()
	s.Events.S2CRegistryEvent.UpdateConnection.AddHandler(func(info message.C2SInfo) { UpdateClient(info.ClientInfo) })
	s.Events.S2CRegistryEvent.UpdateConnection.Enable()
	s.Events.ClientDisconnection.AddHandler(DeleteClient)
	s.Events.ClientDisconnection.Enable()

	S2SEventLogger := func(op, dir string) func(info message.S2SInfo) {
		return func(info message.S2SInfo) {
			log.Println(fmt.Sprintf("%s connection %s %s", op, dir, info.GetServerID()))
		}
	}
	s.Events.S2SRegistrantEvent.NewConnection.AddHandler(S2SEventLogger("New", "to"))
	s.Events.S2SRegistrantEvent.NewConnection.Enable()
	s.Events.S2SRegistrantEvent.UpdateConnection.AddHandler(S2SEventLogger("Update", "to"))
	s.Events.S2SRegistrantEvent.UpdateConnection.Enable()
	s.Events.S2SRegistryEvent.NewConnection.AddHandler(S2SEventLogger("New", "from"))
	s.Events.S2SRegistryEvent.NewConnection.Enable()
	s.Events.S2SRegistryEvent.UpdateConnection.AddHandler(S2SEventLogger("Update", "from"))
	s.Events.S2SRegistryEvent.UpdateConnection.Enable()

	s.Events.S2SRegistrantEvent.Disconnection.AddHandler(func(info message.S2SInfo, err error) {
		log.Println(fmt.Sprintf("Dicsonnection from %s because of %s", info.GetServerID(), err))
	})
	s.Events.S2SRegistrantEvent.Disconnection.Enable()
	s.Events.S2SRegistryEvent.Disconnection.AddHandler(func(info message.S2SInfo) {
		log.Println(fmt.Sprintf("Dicsonnection from %s", info.GetServerID()))
	})
	s.Events.S2SRegistryEvent.Disconnection.Enable()
}
