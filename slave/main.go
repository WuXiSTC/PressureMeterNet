package main

import (
	"PressureMeterNet/slave/option"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	opt, exit := option.Generate(func(i ...interface{}) {
		fmt.Println(i...)
	})
	if exit { //如果要退出
		return //就直接退出
	}
	//fmt.Println(opt)
	mpc := PressureMeterInit(opt.PressureMeterConfig)
	client := ClientInit(opt, mpc)

	ctxBackground := context.Background()
	ctx, cancel := context.WithCancel(ctxBackground)
	errChan := make(chan error, 2)

	go func() {
		client.Run(ctx)
		errChan <- nil
	}()
	if err := <-errChan; err != nil {
		log.Println(err)
		cancel()
	}
	mpc.StopAll()
	time.Sleep(3e9)
}
