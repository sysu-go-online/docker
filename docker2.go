package main

import (
	"log"
	"net"

	"github.com/bilibiliChangKai/First-Docker/socket"
)

var SOCKET_ADDR string = "0.0.0.0:8080" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"

func main() {
	listener, err := net.Listen("tcp", SOCKET_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept() //用conn接收链接
		if err != nil {
			log.Fatal(err)
		}
		go socket.HandleConnection(conn) //开启多个协程。
	}
	// cmd := cmdcreator.Goget("huziang", "github.com/golang/example/hello")
	// runDocker(cmd)
}
