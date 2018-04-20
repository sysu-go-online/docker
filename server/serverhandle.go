package server

//********************************************
// Author : huziang
//   此文件包含路由函数的实现
//********************************************

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/docker/cmdcreator"
	"github.com/sysu-go-online/docker/container"
	"github.com/unrolled/render"
)

const (
	// SocketReadBufferSize Socket套接字的读缓冲区长度
	SocketReadBufferSize = 1024
	// SocketWriteBufferSize Socket套接字的写缓冲区长度
	SocketWriteBufferSize = 1024
)

// HandleConnection 这个是在处理客户端会阻塞的代码。
func HandleConnection(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in?????")
		conn, err := websocket.Upgrade(w, r, nil, SocketReadBufferSize, SocketWriteBufferSize)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// 反json化
		command := &cmdcreator.Command{}
		conn.ReadJSON(command)
		fmt.Println(command)

		// 获得docker命令
		cmd := command.Goget()
		con := container.NewContainer(conn, cmd)

		// container开始运行
		con.Init()

		// 等待container运行结束
		con.Join()
	}
}
