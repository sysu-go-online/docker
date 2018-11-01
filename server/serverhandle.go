package server

//********************************************
// Author : huziang
//   此文件包含路由函数的实现
//********************************************

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/docker_end/container"
	"github.com/sysu-go-online/docker_end/types"
	"github.com/unrolled/render"
)

const (
	// SocketReadBufferSize Socket套接字的读缓冲区长度
	SocketReadBufferSize = 1024
	// SocketWriteBufferSize Socket套接字的写缓冲区长度
	SocketWriteBufferSize = 1024
)

// HandleTTYConnection 这个是在处理客户端会阻塞的代码。
func HandleTTYConnection(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, SocketReadBufferSize, SocketWriteBufferSize)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		msg := types.ConnectContainerRequest{}
		conn.ReadJSON(&msg)

		err = container.StartContainer(msg.ID)
		defer container.StopContainer(msg.ID)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		hijack, err := container.GetHijackRes(msg.ID)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		go container.WriteToUserConn(conn, hijack.Reader)
		_, err = hijack.Conn.Write([]byte(msg.Msg))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		container.WriteToContainer(conn, hijack.Conn)
	}
}

// ContainerCreateHandler create new container and return its id
func ContainerCreateHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}
		msg := types.CreateContainerRequest{}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}
		ID := container.NewContainer(&msg)
		res := types.CreateContainerResponse{}
		if ID == "" {
			res.OK = false
		} else {
			res.OK = true
			res.ID = ID[:12]
		}
		b, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Write(b)
	}
}
