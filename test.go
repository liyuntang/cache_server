package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	for  {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("继续提供服务中...........")
		case do := <-ch:
			fmt.Println("我要关闭了", do)
			fmt.Println("拜拜")
			os.Exit(0)
		}
	}
}