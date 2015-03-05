package util

import (
	"os"
	"path/filepath"
)

// 获得可执行程序所在目录
func ExecutableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}
