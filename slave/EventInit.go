package main

import (
	"PressureMeterNet/slave/option"
	"encoding/json"
	"fmt"
	"github.com/yindaheng98/go-utility/MultiProcessController"
	"github.com/yindaheng98/gogisnet/grpc/client"
	pb "github.com/yindaheng98/gogisnet/grpc/protocol/protobuf"
	"github.com/yindaheng98/gogisnet/message"
	"github.com/yindaheng98/gogistry/protocol"
	"log"
	"strconv"
	"sync"
)

//获取指定数量的IP列表
func GetAddrs(c option.PressureMeterConfig, N uint16) []byte {
	Addrs := make([]string, N)
	for i := uint16(0); i < N; i++ {
		Addrs[i] = c.BoardCastAddr + fmt.Sprintf(":%d", c.BaseListenPort+i)
	}
	addrs, _ := json.Marshal(Addrs)
	return addrs
}
func EventInit(cli *client.Client, mpc *MultiProcessController.MultiProcessController, opt option.PressureMeterConfig) {
	updatePressureMeters := func(si *pb.ServerInfo) { //线程不安全的更新操作
		AdditionalInfo := make(map[string]string, len(si.AdditionalInfo))
		for k, v := range si.AdditionalInfo {
			AdditionalInfo[k] = string(v)
		}
		if TaskAccN, err := strconv.Atoi(string(si.AdditionalInfo["TaskAccN"])); err != nil { //读取新的TaskAccN
			log.Println(fmt.Sprintf("PressureMeter update to %s failed:", si.ServerID), AdditionalInfo, err)
		} else {
			mpc.StartN(uint64(TaskAccN)) //更新MPC
			log.Println(fmt.Sprintf("PressureMeter update to %s succeeded: %s", si.ServerID, AdditionalInfo))
			ci := cli.GetC2SInfo().ClientInfo.(*pb.ClientInfo)
			ci.AdditionalInfo["JmeterAddresses"] = GetAddrs(opt, uint16(TaskAccN)) //更新客户端信息
		}
	}
	pmu := new(sync.Mutex)                            //线程锁
	UpdatePressureMeters := func(si *pb.ServerInfo) { //线程安全的更新操作
		pmu.Lock()
		defer pmu.Unlock()
		updatePressureMeters(si)
	}
	RestartPressureMeters := func(si *pb.ServerInfo) { //线程安全的重启操作
		pmu.Lock()
		defer pmu.Unlock()
		mpc.StopAll()            //先停止
		updatePressureMeters(si) //再启动
	}

	cli.Events.NewConnection.AddHandler(func(info message.S2CInfo) {
		log.Println(fmt.Sprintf("Connected to %s", info.GetServerID()))
		UpdatePressureMeters(info.ServerInfo.(*pb.ServerInfo)) //新连接到来就更新
	})
	cli.Events.NewConnection.Enable()

	cli.Events.Disconnection.AddHandler(func(info message.S2CInfo, err error) {
		log.Println(fmt.Sprintf("Disconnected to %s because of", info.GetServerID()), err)
		RestartPressureMeters(info.ServerInfo.(*pb.ServerInfo)) //断连则重启
	})
	cli.Events.Disconnection.Enable()

	cli.Events.UpdateConnection.AddHandler(func(info message.S2CInfo) {
		log.Println(fmt.Sprintf("Updated connection to %s", info.GetServerID()))
		UpdatePressureMeters(info.ServerInfo.(*pb.ServerInfo)) //更新连接就更新
	})
	cli.Events.UpdateConnection.Enable()

	cli.Events.Retry.AddHandler(func(request protocol.TobeSendRequest, err error) {
		id := request.Request.RegistrantInfo.GetRegistrantID()
		log.Println(fmt.Sprintf("Retried connection to %s because of\n", id), err)
	})
	cli.Events.Retry.Enable()
}
