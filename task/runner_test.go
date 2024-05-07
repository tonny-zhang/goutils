package task

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	dir := ".cache"
	os.MkdirAll(dir, os.ModePerm)
	cDone := make(chan bool)
	var locker sync.Mutex
	lenDone := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano() + 10))
	var runner Runner[int]
	e := runner.Config(Conf[int]{
		DirCache: dir,
		OnRun: func(param int) {
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

	fmt.Println("after add")

	for range time.Tick(1 * time.Second) {
		// fmt.Println("心跳1", time.Now())
	}
	t.Fail()
}
