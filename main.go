package main

import (
	"fmt"
	"io"

	"github.com/bilibiliChangKai/First-Docker/docker"

	"github.com/bilibiliChangKai/First-Docker/cmdcreator"
	"github.com/bilibiliChangKai/First-Docker/socket"
	"github.com/codegangsta/martini"
)

// SOCKETADDR socket地址 string
const (
	socketaddr = ":8081" //表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
)

func mai2() {
	m := martini.Classic()
	m.Get("/", socket.HandleConnection)

	m.RunOnAddr(socketaddr)
}

func main() {
	command := &cmdcreator.Command{}
	command.Command = "go get -v github.com/TXZdream/websocket-echo-server_end"
	cmd := command.Goget()
	pipeReader, _ := io.Pipe()
	outchan := make(chan []byte, 1024)
	errchan := make(chan []byte, 1024)
	go docker.RunDocker(cmd, pipeReader, outchan, errchan)

	// read chan
	outend := false
	errend := false
	for {
		if !outend {
			ob, ok := <-outchan
			if ok == false {
				outend = true
			} else {
				fmt.Println("OUT:", string(ob))
			}
		}

		if !errend {
			ob, ok := <-errchan
			if ok == false {
				errend = true
			} else {
				fmt.Println("ERR:", string(ob))
			}
		}

		if errend && outend {
			break
		}
	}
}
