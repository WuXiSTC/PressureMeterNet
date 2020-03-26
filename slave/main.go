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
	client := ClientInit(opt)
	app := PressureMeterInit(opt.PressureMeterConfig)

	ctxBackground := context.Background()
	ctx, cancel := context.WithCancel(ctxBackground)
	errChan := make(chan error, 2)

	go func() {
		client.Run(ctx)
		errChan <- nil
	}()
	go func() {
		err := app.Start()
		if err != nil {
			errChan <- err
		}
		errChan <- app.Wait()
	}()
	if err := <-errChan; err != nil {
		log.Println(err)
		cancel()
	}
	time.Sleep(3e9)
}
