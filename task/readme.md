## 有最大执行数量的任务队列

## 用法
```go
package main

import (
	"fmt"
	"goutils/task"
	"math/rand"
	"sync"
	"time"
)

var dirCache = "/Users/mac/doc_zk/bitgene/cache"

func main() {
	cDone := make(chan bool)
	var locker sync.Mutex
	lenDone := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano() + 10))
	var runner task.Runner
	e := runner.Config(task.Conf{
		DirCache: dirCache,
		OnRun: func(param interface{}) {
			x := r.Intn(5)
			time.Sleep(time.Duration(x) * time.Second)
			locker.Lock()
			defer locker.Unlock()
			lenDone++
			fmt.Println("外部", param, "已完成", lenDone, x)

			if lenDone == 40 {
				cDone <- true
			}
		},
		OnStop: func() {
			fmt.Println("runner stop")
		},
	})
	fmt.Println("runner 配置", e)
	e = runner.Start()

	if nil == e {
		for i := 0; i < 20; i++ {
			runner.AddTask(i)
		}
	} else {
		fmt.Println(e)
	}

	go func() {
		for i := 0; i < 20; i++ {
			runner.AddTask(i * 100)
		}
		<-cDone
		runner.Stop()
	}()

	for range time.Tick(1 * time.Second) {
		// fmt.Println("心跳1", time.Now())
	}
}

```