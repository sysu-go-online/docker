package main

//********************************************
// Author : huziang
//   启动文件，关联server
//********************************************

import (
	"github.com/sysu-go-online/docker/server"
)

// SOCKETADDR socket地址 string
const (
	addr = ":8081" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
)

func main() {
	m := server.NewServer()
	m.RunOnAddr(addr)
}
