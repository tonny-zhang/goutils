package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/tonny-zhang/goutils/cache"
	"github.com/tonny-zhang/goutils/fileutils"
)

// Conf 配置
type Conf[T any] struct {
	NumRunner int
	DirCache  string
	OnRun     func(param T)
	OnStop    func()
}

// Runner 任务执行实例
type Runner[T any] struct {
	inited      bool
	cRunnerDone chan interface{}
	cRunner     chan byte
	cStop       chan bool
	queue       []T
	isConfed    bool
	isRunning   bool

	numRunner int
	dirCache  string
	onRun     func(param T)
	onStop    func()
}

// Config 配置runner
func (runner *Runner[T]) Config(conf Conf[T]) (err error) {
	_dirCache := conf.DirCache
	if "" == _dirCache {
		err = fmt.Errorf("请先设置缓存目录")
	} else {
		if !fileutils.IsFileExists(_dirCache) {
			err = fmt.Errorf("缓存目录[%s]不存在", _dirCache)
		}
	}
	if nil == err {
		runner.isConfed = true
		runner.dirCache = path.Join(conf.DirCache, ".task_wait")
		runner.onRun = conf.OnRun
		runner.onStop = conf.OnStop
		runner.numRunner = conf.NumRunner
	}
	return
}

func (runner *Runner[T]) check() (err error) {
	if !runner.isConfed {
		err = fmt.Errorf("请先调用Config方法")
	}

	return
}

// Start 开始执行
func (runner *Runner[T]) Start() (err error) {
	err = runner.check()
	if err == nil {
		if !runner.isRunning {
			if runner.numRunner == 0 {
				runner.numRunner = runtime.NumCPU()
			}
			runner.cRunnerDone = make(chan interface{}, runner.numRunner)
			runner.cRunner = make(chan byte, runner.numRunner)
			runner.cStop = make(chan bool)

			fileInfoList, err := os.ReadDir(runner.dirCache)

			if nil == err {
				// fmt.Printf("有%d个旧任务需要处理\n", len(fileInfoList))
				for _, file := range fileInfoList {
					filepath := path.Join(runner.dirCache, file.Name())

					content, _ := cache.GetCache(filepath, 0)
					var param T
					e := json.Unmarshal(content, &param)
					if nil == e {
						runner.AddTask(param)
					}
				}
			}
			go func() {
			loop:
				for {
					select {
					case <-runner.cRunnerDone:
						if len(runner.queue) > 0 {
							param := runner.queue[0]
							runner.queue = runner.queue[1:]

							runner.runTask(param)
						}
					case <-runner.cStop:
						runner.isRunning = false
						break loop
					}
				}
				if nil != runner.onStop {
					runner.onStop()
				}
			}()
			runner.isRunning = true
		}
	}
	return
}

// Stop 停止任务
func (runner *Runner[T]) Stop() {
	runner.cStop <- true
}

func (runner *Runner[T]) runTask(paramTask T) {
	runner.cRunner <- 0
	go func() {
		if nil != runner.onRun {
			runner.onRun(paramTask)
		}

		// key := getMD5(paramTask)
		// filecache := path.Join(runner.dirCache, key)
		// cache.RemoveCache(filecache)
		<-runner.cRunner
		runner.cRunnerDone <- paramTask
	}()
}

// AddTask 添加任务
func (runner *Runner[T]) AddTask(paramTask T) (err error) {
	err = runner.check()
	if nil == err {
		key := getMD5(paramTask)
		filecache := path.Join(runner.dirCache, key)

		if !fileutils.IsFileExists(filecache) {
			cache.SetCacheJSON(filecache, paramTask)
		}

		if len(runner.cRunner) < runner.numRunner {
			runner.runTask(paramTask)
		} else {
			runner.queue = append(runner.queue, paramTask)
		}
	}
	return
}
