package env

import (
	"path"
	"sync"

	"github.com/tonny-zhang/goutils/fileutils"

	"github.com/joho/godotenv"
)

var locker sync.Once

// AutoLoad 自动加载环境变量
func AutoLoad() {
	locker.Do(func() {
		// 自动加载环境变量
		var envFiles []string
		envFile := path.Join(fileutils.GetCmdDir(), ".env")
		if fileutils.IsFileExists(envFile) {
			envFiles = append(envFiles, envFile)
		}
		godotenv.Overload(envFiles...)
	})
}
