package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 获取可执行文件所在的目录
func GetExecuteFileDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = path[:index]
	return path
}
