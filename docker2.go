package main

import (
	"github.com/codegangsta/martini"

	"github.com/bilibiliChangKai/First-Docker/socket"
)

// SOCKETADDR socket地址 string
const (
	socketaddr = ":8081" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
)

func main() {
	m := martini.Classic()
	m.Get("/", socket.HandleConnection)

	m.RunOnAddr(socketaddr)
}
