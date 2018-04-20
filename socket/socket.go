package socket

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/docker/cmdcreator"
	"github.com/sysu-go-online/docker/docker"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// Signal 异步信号量
type Signal struct {
	num int
	ch  chan bool
}

// NewSignal 新创建一个信号量指针
func NewSignal(num int) *Signal {
	signal := new(Signal)
	signal.num = num
	return signal
}

// Signal 发送信号
func (signal *Signal) Signal() {
	signal.num--
	signal.ch <- true
	return
}

// Wait 等待信号
func (signal *Signal) Wait() {
	for signal.num > 0 {
		<-signal.ch
	}
}

// 异步读取
func asycReadChannel(conn *websocket.Conn, ch chan []byte, signal *Signal) {
	for {
		ob, ok := <-ch
		if ok == false {
			signal.Signal()
			break
		} else {
			conn.WriteMessage(websocket.TextMessage, ob)
		}
	}
}

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
	conn.ReadJSON(command)
	fmt.Println(command)

	conn.WriteMessage(websocket.TextMessage, []byte("Welcome connect!"))

	signal := NewSignal(2)
	outchan := make(chan []byte, 1024)
	errchan := make(chan []byte, 1024)

	// 运行docker命令
	cmd := command.Goget()
	pipeReader, _ := io.Pipe()
	go docker.RunDocker(cmd, pipeReader, outchan, errchan)

	// asyc read two signal
	go asycReadChannel(conn, outchan, signal)
	go asycReadChannel(conn, errchan, signal)

	signal.Wait()
}
