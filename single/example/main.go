package main

import (
	"fmt"
	"goutils/single"
)

func main() {
	// 自动使用sock文件
	// started := single.StartAuto()

	// 指定sock文件
	started := single.Start("./1.sock")
	if started {
		fmt.Println("已经启动了一个实例")
	} else {
		fmt.Println("第一次运行")
		select {}
	}
}
