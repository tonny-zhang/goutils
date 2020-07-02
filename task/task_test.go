package task

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

var fCurrent, _ = filepath.Abs(os.Args[0])
var dirCurrent = path.Dir(fCurrent)

var dirCache = "/Users/mac/doc_zk/bitgene/cache"

func TestAddTaskNoConf(t *testing.T) {
	err := AddTask(nil)
	fmt.Println(err)
}

func TestConf(t *testing.T) {
	err := SetConf(Conf{
		DirCache: dirCache,
	})
	fmt.Println(err)

	AddTask(123)

	// select {}
	time.Sleep(100 * time.Second)
}
