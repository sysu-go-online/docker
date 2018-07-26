package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//********************************************
// Author : huziang
//   包含常用函数
//********************************************

// DealPanic 处理错误
func DealPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func isFileExit() bool {
	return true
}

// GetGOPATH 获得用户环境的gopath
func GetGOPATH() string {
	var sp string
	if runtime.GOOS == "windows" {
		sp = ";"
	} else {
		sp = ":"
	}
	goPath := strings.Split(os.Getenv("GOPATH"), sp)
	for _, v := range goPath {
		if _, err := os.Stat(filepath.Join(v, "/src/github.com/sysu-go-online/docker_end/util/util.go")); !os.IsNotExist(err) {
			return v
		}
	}
	return ""
}
