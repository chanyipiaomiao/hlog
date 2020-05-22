package hlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// isExist 文件或目录是否存在
// return false 表示文件不存在
func isExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// 创建日志目录
func makeDirAll(logPath string) error {

	logDir := path.Dir(logPath)
	if !isExist(logDir) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("create <%s> error: %s", logDir, err)
		}
	}

	return nil
}

// 是否是windows系统
func isWindow() bool {
	return runtime.GOOS == "windows"
}
