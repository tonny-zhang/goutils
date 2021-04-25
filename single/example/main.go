package main

import (
	"fmt"
	"goutils/single"
)

func main() {
	started := single.StartAuto()

	if started {
		fmt.Println("已经启动了一个实例")
	} else {
		fmt.Println("第一次运行")
		select {}
	}
}
