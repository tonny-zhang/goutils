package env

import (
	"path"
	"path/filepath"
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
			if envFileCurrentPath, e := filepath.Abs(".env"); e == nil {
				if envFileCurrentPath != envFile && fileutils.IsFileExists(envFileCurrentPath) {
					envFiles = append(envFiles, envFileCurrentPath)
				}
			}
		} else {
			envFiles = append(envFiles, ".env")
		}

		godotenv.Overload(envFiles...)
	})
}

// LoadFromFile 从指定文件加载
func LoadFromFile(envFiles ...string) {
	godotenv.Overload(envFiles...)
}
