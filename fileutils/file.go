package fileutils

import "os"

// IsFileExists 查看指定文件是否存在
func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Mkdirp 目录没有时创建
func Mkdirp(dir string) {
	if !IsFileExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}
