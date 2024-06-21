package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/tonny-zhang/goutils/fileutils"
)

var defaultWriter io.Writer = os.Stdout

// Logger logger
type Logger struct {
	prefix          string
	writer          io.Writer
	printStack      bool // 是否打印发生错误路径
	CloseLog        bool
	HideProjectPath bool // 是否隐藏项目路径
}

// SetWriter set writer for logger
func (logger *Logger) SetWriter(writer io.Writer) {
	logger.writer = writer
}

// SetWriter 设置默认输出对象
func SetWriter(writer io.Writer) {
	defaultWriter = writer

	defaultLogger = Logger{
		writer: defaultWriter,
	}
}

// Info info for log
func (logger Logger) log(prev, formater string, msg ...any) {
	if logger.CloseLog {
		return
	}
	writer := logger.writer
	if writer == nil {
		writer = defaultWriter
	}
	formater = "%s %-8s %s \t" + formater
	msg = append([]any{
		"[" + time.Now().Format("2006/01/02 15:04:05") + "]",
		"[" + prev + "]",
		logger.prefix,
	}, msg...)
	msgToWrite := fmt.Sprintf(formater, msg...)
	if _, e := fmt.Fprintln(writer, msgToWrite); e != nil {
		// write log to stdout
		fmt.Fprintln(os.Stdout, msgToWrite)
	}
}

// Info info for log
func (logger Logger) Info(formater string, msg ...any) {
	logger.log("Info", formater, msg...)
}

// Warn warn for log
func (logger Logger) Warn(formater string, msg ...any) {
	logger.log("Warn", formater, msg...)
}

func (logger Logger) Error(formater string, msg ...any) {
	if logger.printStack {
		fileNumInfo := ""
		_, filePath, line, _ := runtime.Caller(1)

		if logger.HideProjectPath {
			basedir := fileutils.GetRuntimeProjectBaseDir()
			filePath = strings.Replace(filePath, basedir, "", -1)
		}

		fileNumInfo = fmt.Sprintf("[%s:%d]", filePath, line)

		formater = "%s : " + formater
		msg = append([]any{fileNumInfo}, msg...)
	}

	logger.log("Error", formater, msg...)
}

// Debug debug for log
func (logger Logger) Debug(formater string, msg ...any) {
	logger.log("Debug", formater, msg...)
}

var defaultLogger = Logger{
	writer:          defaultWriter,
	HideProjectPath: true,
	printStack:      true,
}
var loggerMap = make(map[string]Logger)

// DefaultLogger get default logger
func DefaultLogger() Logger {
	return defaultLogger
}

// PrefixLogger 得到有前缀输出的logger
func PrefixLogger(prefix string) Logger {
	if logger, ok := loggerMap[prefix]; ok {
		return logger
	}
	logger := Logger{
		prefix:          prefix,
		HideProjectPath: true,
		printStack:      true,
	}

	loggerMap[prefix] = logger

	return logger
}
