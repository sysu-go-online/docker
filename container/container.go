package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/gorilla/websocket"
)

const (
	// BufferSize 缓冲区大小，以字节为单位
	BufferSize = 1024
)

var idset int

// Container 通过接口封装输入输出给
type Container struct {
	ID    int             // 主键
	conn  *websocket.Conn // 绑定的websocket，其中一端
	cmd   *exec.Cmd       // 绑定的cmd执行指令，另外一段
	isEnd *Signal         // 判断是否结束
}

// NewContainer 新创建一个容器指针
func NewContainer(conn *websocket.Conn, cmd *exec.Cmd) *Container {
	container := new(Container)
	container.ID = idset
	idset++
	container.conn = conn
	container.cmd = cmd
	container.isEnd = NewSignal(2)
	return container
}

// Init 协程初始化，并开始运行
func (con *Container) Init() {
	// docker stdin, stdout, stderr设置
	//con.cmd.Stdin = bufio.NewReader(in)
	stdout, err := con.cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := con.cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	stdin, err := con.cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	// docker 启动
	con.cmd.Start()

	// 异步读取信息，并发送给connection
	asynWriteChannel := func(out io.ReadCloser) {
		defer out.Close()
		reader := bufio.NewReader(out)
		for {
			line, _, err := reader.ReadLine()
			// 传输结束，发出信号
			if err != nil {
				con.isEnd.Signal()
				break
			} else {
				con.conn.WriteMessage(websocket.TextMessage, line)
			}
		}
	}

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
			_, err = in.Write(msg)
			// If message can not be written to the process, kill it
			if err != nil {
				con.cmd.Process.Kill()
				break
			}
		}
	}

	// asyn write two channel
	go readFromConnection(stdin)
	go asynWriteChannel(stdout)
	go asynWriteChannel(stderr)
}

// Join 当两个协程全部运行完毕时，程序结束
func (con *Container) Join() {
	con.isEnd.Wait()
}
