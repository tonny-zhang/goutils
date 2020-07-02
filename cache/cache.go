package cache

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/tonny-zhang/goutils/fileutils"
)

// // IsFileExists 查看指定文件是否存在
// func isFileExists(name string) bool {
// 	if _, err := os.Stat(name); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }

// // Mkdirp 目录没有时创建
// func mkdirp(dir string) {
// 	if !isFileExists(dir) {
// 		os.MkdirAll(dir, os.ModePerm)
// 	}
// }

// RemoveCache 删除缓存文件
func RemoveCache(filecache string) {
	os.Remove(filecache)
}

// GetCache 得到缓存
func GetCache(filecache string, delayHours int64) (b []byte, err error) {
	if fileutils.IsFileExists(filecache) {
		if delayHours > 0 {
			info, _ := os.Stat(filecache)
			m, _ := time.ParseDuration(strconv.FormatInt(delayHours, 10) + "h")
			timeDelay := info.ModTime().Add(m)
			if time.Now().After(timeDelay) {
				err = fmt.Errorf("过期")
				return
			}
		}
		b, err = ioutil.ReadFile(filecache)
	} else {
		err = fmt.Errorf("不存在")
	}
	return
}

// SetCacheJSON 设置缓存，内容为json
func SetCacheJSON(filecache string, data interface{}) {
	buf, e := json.Marshal(data)
	if e == nil {
		SetCache(filecache, buf)
	}
}

// SetCache 设置缓存
func SetCache(filecache string, content []byte) {
	fileutils.Mkdirp(path.Dir((filecache)))
	fw, e := os.Create(filecache)
	if e == nil {
		defer fw.Close()
		fw.Write(content)
	}
}
