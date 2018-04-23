package main

//********************************************
// Author : huziang
//   启动文件，关联server
//********************************************

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/sysu-go-online/docker_end/server"
)

// SOCKETADDR socket地址 string
// var (
// 	addr = ":8081" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
// )

func main() {
	// Get port from env
	var PORT = os.Getenv("PORT")
	if len(PORT) == 0 {
		PORT = "8888"
	}
	var port = pflag.StringP("port", "p", PORT, "Define the port where service runs")
	pflag.Parse()
	addr := ":" + *port

	m := server.NewServer()
	m.RunOnAddr(addr)
}
