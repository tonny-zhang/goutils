package fileutils

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

var dirRuntimeProjectBase = "" // 运行时项目根目录
var initOnce sync.Once

var isRunTypeSource = false // 是否从源码运行
var cmdBaseDir string       // 当前可执行二进制文件运行根目录

// IsFileExists 查看文件是否存在
func IsFileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsLinkExists 软链接是否存在
func IsLinkExists(name string) bool {
	if _, err := os.Lstat(name); err != nil {
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

var regTrim = regexp.MustCompile("\\s*")

func parseRuntimeProjectBaseDir(stack, mainModuleName string) (res string) {
	if mainModuleName != "" {
		/**

		15 main.main()
		16      {{项目路径}}/main.go:30 +0x6c
		------------------
		13 main.main.func1()
		14      {{项目路径}}/main.go:33 +0x6b
		15 created by main.main in goroutine 1
		16      {{项目路径}}/main.go:32 +0x1a
		------------------
		13 pcai/lib.Run.func1()
		14      {{项目路径}}/lib/test.go:25 +0x6b
		15 created by pcai/lib.Run in goroutine 1
		16      {{项目路径}}/lib/test.go:24 +0x76
		*/
		arr := strings.Split(stack, "\n")
		for i := len(arr) - 2; i > 0; i-- {
			if strings.HasPrefix(arr[i], "main.main") {
				res = strings.Replace(strings.Split(regTrim.ReplaceAllString(arr[i+1], ""), ":")[0], "/main.go", "", -1)
				break
			} else if strings.HasPrefix(arr[i], mainModuleName) {
				pathInfo := arr[i+1]
				moduleInfo := strings.Replace(arr[i], mainModuleName, "", -1)

				indexLastSplit := strings.LastIndex(moduleInfo, "/")
				moduleInfo2 := moduleInfo[0:indexLastSplit+1] + strings.Split(moduleInfo[indexLastSplit+1:], ".")[0]

				moduleInfo2 = strings.ReplaceAll(moduleInfo2, "%2e", ".")

				index := strings.LastIndex(pathInfo, "/")
				if index > -1 {
					pathInfo = pathInfo[0:index]
					res = pathInfo[0:strings.LastIndex(pathInfo, moduleInfo2)]
					res = regTrim.ReplaceAllString(res, "")

					break
				}
			}
		}
	}

	return
}
func initData() {
	initOnce.Do(func() {
		if executable, err := os.Executable(); err == nil {
			// fmt.Println("executable", executable)
			switch runtime.GOOS {
			case "windows":
				// C:\Users\5950X\AppData\Local\Temp\GoLand\___go_build_github_com_golang_infrastructure_go_project_root_directory_main_test.exe
				if strings.Contains(executable, "\\AppData\\Local\\Temp\\") {
					isRunTypeSource = true
				}
			case "linux":
				// /tmp/go-build1325605723/b001/exe/test
				if strings.HasPrefix(executable, "/tmp/go-build") {
					isRunTypeSource = true
				}
			case "darwin":
				// /var/folders/kd/dzyx8fc96fx4j3mtdtjsl4z40000gn/T/go-build3362823274/b001/exe/main
				if strings.Contains(executable, "/T/go-build") {
					isRunTypeSource = true
				}
			}
		}

		if isRunTypeSource {
			if searchDirectory, err := os.Getwd(); err == nil {
				// 从当前路径往上找，第一个拥有go.mod文件的目录认为是项目的根路径
				for searchDirectory != "" {
					goModPath := filepath.Join(searchDirectory, "go.mod")
					stat, err := os.Stat(goModPath)
					if err == nil && stat != nil && !stat.IsDir() {
						cmdBaseDir = searchDirectory
						dirRuntimeProjectBase = searchDirectory
						break
					}
					searchDirectory = filepath.Dir(searchDirectory)
				}
			}
		} else {
			if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
				cmdBaseDir = dir
			}

			if v, ok := debug.ReadBuildInfo(); ok {
				mainModuleName := v.Main.Path
				if mainModuleName == "" {
					if len(v.Deps) > 0 {
						mainModuleName = v.Deps[len(v.Deps)-1].Path
					}
				}

				if v := parseRuntimeProjectBaseDir(string(debug.Stack()), mainModuleName); v != "" {
					dirRuntimeProjectBase = v
				}
			}
		}
	})
}

// GetRuntimeProjectBaseDir 得到当前主项目的根目录
func GetRuntimeProjectBaseDir() string {
	initData()
	return dirRuntimeProjectBase
}

// GetCmdDir 得到当前可执行文件所在根目录，调试时得到源码项目根目录
func GetCmdDir() string {
	initData()
	return cmdBaseDir
}
