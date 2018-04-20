package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sysu-go-online/docker/docker"

	"github.com/codegangsta/martini"
	"github.com/sysu-go-online/docker/socket"
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

func main2() {
	// command := &cmdcreator.Command{}
	// command.Command = "go get -v github.com/TXZdream/websocket-echo-server_end"
	// cmd := command.Goget()
	cmdName := "go get -v github.com/TXZdream/websocket-echo-server"
	cmdArgs := strings.Fields(cmdName)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	pipeReader, _ := io.Pipe()
	outchan := make(chan []byte, 1024)
	errchan := make(chan []byte, 1024)
	go docker.RunDocker(cmd, pipeReader, outchan, errchan)

	// read chan
	outend := false
	errend := false
	go func() {
		for {
			ob, ok := <-outchan
			if ok == false {
				outend = true
				break
			} else {
				fmt.Println("OUT:", string(ob))
			}
		}
	}()

	go func() {
		for {
			ob, ok := <-errchan
			if ok == false {
				errend = true
				break
			} else {
				fmt.Println("ERR:", string(ob))
			}
		}
	}()

	for {
		if outend && errend {
			break
		}
	}

	// for {
	// 	if !outend {
	// 		ob, ok := <-outchan
	// 		if ok == false {
	// 			outend = true
	// 		} else {
	// 			fmt.Println("OUT:", string(ob))
	// 		}
	// 	}

	// 	if !errend {
	// 		ob, ok := <-errchan
	// 		if ok == false {
	// 			errend = true
	// 		} else {
	// 			fmt.Println("ERR:", string(ob))
	// 		}
	// 	}

	// 	if errend && outend {
	// 		break
	// 	}
	// }
}
