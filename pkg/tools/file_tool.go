package tools

import (
	"flygoose/pkg/tlog"
	"os"
	"path/filepath"
)

// 获取可执行文件的目录
func GetExecuteDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	exPath := filepath.Dir(ex)
	return exPath
}

// 判断文件或者文件夹是否存在
func IsFileExist(path string) bool {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false
	}
	if err == nil {
		return true
	}
	return false
}

func CreateDir(targetDir string) string {
	if !IsFileExist(targetDir) {
		err := os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			tlog.Error2("创建目录失败", err)
			return ""
		} else {
			return targetDir
		}
	}
	return targetDir
}
