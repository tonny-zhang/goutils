package task

import (
	"fmt"
	"path"
	"runtime"

	"cache"
	"fileutils"
)

// 是否已经配置过
var _isSetedConf bool

// 用于传递任务参数
var channelRunner chan interface{}

// 任务完成时调用外部方法
var _onTaskDone func(param interface{})

// 缓存目录
var _dirCache string

var numRunner int
var _addedTask = false

// Conf 配置
type Conf struct {
	NumRunner  int
	DirCache   string
	OnTaskDone func(param interface{})
}

// SetConf 外部配置
func SetConf(conf Conf) (err error) {
	if !_isSetedConf {
		_onTaskDone = conf.OnTaskDone
		numRunner := conf.NumRunner
		if numRunner == 0 {
			numRunner = runtime.NumCPU()
		}
		channelRunner = make(chan interface{}, numRunner)

		_dirCache = conf.DirCache
		if "" == _dirCache {
			err = fmt.Errorf("请先设置缓存目录")
		} else {
			if !fileutils.IsFileExists(_dirCache) {
				err = fmt.Errorf("缓存目录[%s]不存在", _dirCache)
			}
		}
	} else {
		err = fmt.Errorf("已经被初始")
	}
	if err == nil {
		_isSetedConf = true

		fmt.Println("开始监听")
		go start()
		fmt.Println("开始监听2")
	}
	return
}

func _check() (err error) {
	if !_isSetedConf {
		err = fmt.Errorf("请先配置")
	}
	return
}

func start() {
	select {
	case e, param := <-channelRunner:
		fmt.Println("结果", e, param)
	}

}
func runTask(paramTask interface{}) {
	channelRunner <- paramTask
}

// AddTask 添加任务体
func AddTask(paramTask interface{}) (err error) {
	err = _check()
	if nil == err {
		key := getMD5(paramTask)
		filecache := path.Join(_dirCache, ".task_dealing/"+key)

		if !fileutils.IsFileExists(filecache) {
			cache.SetCacheJSON(filecache, paramTask)
		}

		fmt.Println(filecache)

		go runTask(paramTask)
	}
	return
}
