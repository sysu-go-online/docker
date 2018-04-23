package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/gorilla/websocket"

	. "github.com/sysu-go-online/docker_end/util"
)

const (
	// BufferSize 缓冲区大小，以字节为单位
	BufferSize = 1024
)

var idset int

// Container 通过接口封装输入输出给
type Container struct {
	ID   int             // 主键
	conn *websocket.Conn // 绑定的websocket，其中一端
	cmd  *exec.Cmd       // 绑定的cmd执行指令，另外一段
}

// NewContainer 新创建一个容器指针
func NewContainer(conn *websocket.Conn, cmd *exec.Cmd) *Container {
	container := new(Container)
	container.ID = idset
	idset++
	container.conn = conn
	container.cmd = cmd
	return container
}

// Init 协程初始化，并开始运行
func (con *Container) Init() {
	defer func() {
		if err := recover(); err != nil {
			con.conn.WriteMessage(websocket.TextMessage, []byte("Error!"))
		}
	}()

	// docker stdin, stdout, stderr设置
	stdin, err := con.cmd.StdinPipe()
	DealPanic(err)

	stdout, err := con.cmd.StdoutPipe()
	DealPanic(err)

	stderr, err := con.cmd.StderrPipe()
	DealPanic(err)

	// docker 启动
	con.cmd.Start()

	// 异步读取信息，并发送给connection
	writeToConnection := func(out io.ReadCloser) {
		defer out.Close()
		reader := bufio.NewReader(out)
		for {
			line, _, err := reader.ReadLine()
			fmt.Println("Message send: " + string(line))
			// 传输结束，发出信号
			if err != nil {
				break
			} else {
				con.conn.WriteMessage(websocket.TextMessage, line)
			}
		}
	}

	// 异步读取信息，并发送给cmd
	readFromConnection := func(in io.WriteCloser) {
		defer in.Close()

		// Read message from client and write to process
		for {
			_, msg, err := con.conn.ReadMessage()
			// If client close connection, kill the process
			if err != nil {
				con.cmd.Process.Kill()
				break
			}
			fmt.Println("Message get: " + string(msg))
			msg = append(msg, byte('\n'))
			_, err = in.Write(msg)
			// If message can not be written to the process, kill it
			if err != nil {
				con.cmd.Process.Kill()
				break
			}
		}
	}

	// asyn write and read
	go readFromConnection(stdin)
	go writeToConnection(stdout)
	go writeToConnection(stderr)
}

// Join 当两个协程全部运行完毕时，程序结束
func (con *Container) Join() {
	con.cmd.Wait()
}
