package main

import (
	"PressureMeterNet/slave/option"
	"fmt"
	"github.com/yindaheng98/gogisnet/grpc/client"
	"github.com/yindaheng98/gogisnet/message"
	"github.com/yindaheng98/gogistry/protocol"
	"log"
	"net"
	"os/exec"
)

func ClientInit(opt option.Option) *client.Client {
	BoardCastAddr := opt.PressureMeterConfig.BoardCastAddr
	if TCPaddr, err := net.ResolveTCPAddr("", BoardCastAddr); err == nil {
		BoardCastAddr = TCPaddr.String()
	} else if IPaddr, err := net.ResolveIPAddr("", BoardCastAddr); err == nil {
		BoardCastAddr = (&net.TCPAddr{
			IP:   IPaddr.IP,
			Port: opt.PressureMeterConfig.ListenPort,
			Zone: IPaddr.Zone,
		}).String()
	} else {
		panic(err)
	}
	opt.PressureMeterConfig.BoardCastAddr = BoardCastAddr
	opt.ClientInfoOption.AdditionalInfo = BoardCastAddr
	c := client.New(opt.ClientInfoOption, opt.GogisnetOption)
	EventInit(c)
	return c
}

func EventInit(cli *client.Client) {
	cli.Events.NewConnection.AddHandler(func(info message.S2CInfo) {
		log.Println(fmt.Sprintf("Connected to %s", info.GetServerID()))
	})
	cli.Events.NewConnection.Enable()
	cli.Events.Disconnection.AddHandler(func(info message.S2CInfo, err error) {
		log.Println(fmt.Sprintf("Disconnected to %s because of\n", info.GetServerID()), err)
	})
	cli.Events.Disconnection.Enable()
	cli.Events.UpdateConnection.AddHandler(func(info message.S2CInfo) {
		log.Println(fmt.Sprintf("Updated connection to %s", info.GetServerID()))
	})
	cli.Events.UpdateConnection.Enable()
	cli.Events.Retry.AddHandler(func(request protocol.TobeSendRequest, err error) {
		id := request.Request.RegistrantInfo.GetRegistrantID()
		log.Println(fmt.Sprintf("Retried connection to %s because of\n", id), err)
	})
	cli.Events.Retry.Enable()
}

func PressureMeterInit(opt option.PressureMeterConfig) *exec.Cmd {
	app := exec.Command("jmeter",
		fmt.Sprintf("-Jserver.rmi.localport=%d", opt.ListenPort),
		fmt.Sprintf("-Dserver_port=%d", opt.ListenPort),
		"--server", "-Jserver.rmi.ssl.disable=true")
	stdout, err := app.StdoutPipe()
	app.Stderr = app.Stdout
	if err != nil {
		panic(err)
	}
	go func() {
		tmp := make([]byte, 1024*8)
		for err = nil; err == nil; _, err = stdout.Read(tmp) {
			log.Print(string(tmp))
		}
	}()
	return app
}
