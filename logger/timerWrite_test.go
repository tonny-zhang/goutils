package logger

import (
	"fmt"
	"sync"
	"testing"
)

func init() {
	// p, _ := filepath.Abs("./.cache/")
	writer := &TimeWriter{
		LogDir:         "./.cache/",
		FilepathFormat: "20060102.log",
		LastFileName:   "last.log",
	}
	writer.KeepDays = 3

	SetWriter(writer)

}

func TestTimerWrite(t *testing.T) {
	log := DefaultLogger()

	num := 100
	group := sync.WaitGroup{}

	group.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			log.Info("test %d", i)
			group.Done()
		}()
	}

	group.Wait()
	fmt.Println("down")

}
