package socket

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bilibiliChangKai/First-Docker/cmdcreator"
	"github.com/bilibiliChangKai/First-Docker/docker"
	"github.com/gorilla/websocket"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// HandleConnection 这个是在处理客户端会阻塞的代码。
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in?????")
	conn, err := websocket.Upgrade(w, r, nil, readBufferSize, writeBufferSize)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 反json化
	command := &cmdcreator.Command{}

	// _, bs, err := conn.ReadMessage()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(bs))
	conn.ReadJSON(command)
	fmt.Println(command)

	conn.WriteMessage(websocket.TextMessage, []byte("Welcome connect!"))

	outchan := make(chan []byte, 1024)
	errchan := make(chan []byte, 1024)

	// 运行docker命令
	cmd := command.Goget()
	pipeReader, _ := io.Pipe()
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
				conn.WriteMessage(websocket.TextMessage, ob)
			}
		}

		if !errend {
			ob, ok := <-errchan
			fmt.Println("OB:" + string(ob))
			if ok == false {
				errend = true
			} else {
				fmt.Println("ERR:", string(ob))
				conn.WriteMessage(websocket.TextMessage, ob)
			}
		}

		if errend && outend {
			break
		}
	}
}
