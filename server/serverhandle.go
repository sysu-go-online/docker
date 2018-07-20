package server

//********************************************
// Author : huziang
//   此文件包含路由函数的实现
//********************************************

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/docker_end/cmdcreator"
	"github.com/sysu-go-online/docker_end/container"
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
		conn, err := websocket.Upgrade(w, r, nil, SocketReadBufferSize, SocketWriteBufferSize)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// 反json化
		command := &cmdcreator.Command{}
		conn.ReadJSON(command)

		// 新建容器
		con := container.NewContainer(conn, command)
		container.StartContainer(con)
	}
}

// HandleAuth 处理权限验证
// func HandleAuth(formatter *render.Render) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		r.ParseForm()
// 		code := r.Form["code"][0]

// 		tk := ""
// 		token.New(code, &tk)

// 		w.Header().Set("Token", tk)

// 		formatter.JSON(w, 200, struct {
// 			Name string `json:"name"`
// 			Icon string `json:"icon"`
// 		}{})
// 	}
// }
