package logger

import (
	"os"
	"path"
	"sync"
	"time"
)

// TimeWriter 时间wirter
type TimeWriter struct {
	LogDir         string // 日志目录
	FilepathFormat string // 文件名格式
	LastFileName   string // 最新日志快捷文件
	KeepDays       int    // 保留天数

	filepath string
	file     *os.File

	cleaning    bool
	cleanLocker sync.Mutex

	createLinkLocker sync.Mutex // 创建快捷文件锁
}

var loggerForWriter = PrefixLogger("[logger]")

func init() {
	loggerForWriter.PrintStack = false
}

func (writer *TimeWriter) clean() {
	writer.cleanLocker.Lock()
	defer writer.cleanLocker.Unlock()
	defer func() {
		writer.cleaning = false
	}()
	loggerForWriter.Debug("开始清理日志, KeepDays = %v, cleaning = %v", writer.KeepDays, writer.cleaning)

	if writer.KeepDays <= 0 {
		return
	}

	if writer.cleaning {
		return
	}
	writer.cleaning = true

	if list, e := os.ReadDir(writer.LogDir); e == nil {
		t := time.Now().Add(-time.Hour * 24 * time.Duration(writer.KeepDays))

		for _, f := range list {
			logFile := path.Join(writer.LogDir, f.Name())
			if f.IsDir() {
				if list2, e2 := os.ReadDir(logFile); e2 == nil {
					if len(list2) == 0 {
						if e := os.RemoveAll(logFile); e == nil {
							loggerForWriter.Info("删除空目录 [%s]", logFile)
						}
					} else {
						for _, f2 := range list2 {
							if !f2.IsDir() {
								logFile2 := path.Join(logFile, f.Name())
								if info, ef2 := os.Stat(logFile2); ef2 == nil && info.ModTime().Before(t) {
									if e := os.RemoveAll(logFile2); e == nil {
										loggerForWriter.Info("删除日志文件 [%s]", logFile2)
									}
								}
							}
						}
					}
				}
			} else {
				if info, ef2 := os.Stat(logFile); ef2 == nil {
					if info.ModTime().Before(t) {
						if e := os.RemoveAll(logFile); e == nil {
							loggerForWriter.Info("删除日志文件 [%s]", logFile)
						}
					}
				}
			}
		}
	} else {
		loggerForWriter.Error("删除日志文件 [%v]", e)
	}
}
func (writer *TimeWriter) Write(p []byte) (n int, err error) {
	logfilePath := path.Join(writer.LogDir, time.Now().Format(writer.FilepathFormat))
	if writer.filepath != logfilePath && writer.file != nil {
		writer.file.Close()
		writer.file = nil
		writer.filepath = ""

		go writer.clean() // 定时清除日志
	}
	if writer.file == nil {
		writer.createLinkLocker.Lock()
		defer writer.createLinkLocker.Unlock()

		dir := path.Dir(logfilePath)
		e := os.MkdirAll(dir, os.ModePerm)
		if e != nil {
			err = e
			loggerForWriter.Error("创建目录[%s] 异常 [%v]", dir, e)
			return
		}
		writer.filepath = logfilePath
		writer.file, e = os.OpenFile(logfilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if e != nil {
			err = e
			loggerForWriter.Error("创建日志文件[%s] 异常 [%v]", logfilePath, e)
			return
		} else if writer.LastFileName != "" {
			fileLink := path.Join(dir, writer.LastFileName)

			toCreateLink := true
			if target, e := os.Readlink(fileLink); e == nil {
				if target != logfilePath {
					if e := os.Remove(fileLink); e != nil {
						loggerForWriter.Error("删除快捷方式[%s] 异常 [%v]", fileLink, e)
					}
				} else {
					toCreateLink = false // 快捷方式正确指向时不用创建
				}
			}

			if toCreateLink {
				if e := os.Symlink(logfilePath, fileLink); e != nil {
					loggerForWriter.Error("创建最新日志文件快捷方式[%s] 异常 [%v]", fileLink, e)
				} else {
					loggerForWriter.Info("创建最新日志文件快捷方式 %s => %s", logfilePath, fileLink)
				}
			}
		}
	}

	// NOTICE: write after file removed gives error
	// https://stackoverflow.com/questions/34325128/write-to-non-existing-file-gives-no-error/34325329
	n, err = writer.file.Write(p)

	if err != nil {
		loggerForWriter.Error("写入数据[%s] 异常 [%v]", string(p), err)
	}
	return
}
