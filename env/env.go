package env

import (
	"github.com/tonny-zhang/goutils/fileutils"
	"path"

	"github.com/joho/godotenv"
)

func init() {
	// 自动加载环境变量
	var envFiles []string
	envFile := path.Join(fileutils.GetCmdDir(), ".env")
	if fileutils.IsFileExists(envFile) {
		envFiles = append(envFiles, envFile)
	}
	godotenv.Overload(envFiles...)
}
