package socket

import (
	"bufio"
	"io"
	"net"
	"strings"

	"github.com/bilibiliChangKai/First-Docker/cmdcreator"
	"github.com/bilibiliChangKai/First-Docker/docker"
)

// undefined
func getCommand(line string) map[string]string {
	return nil
}

//这个是在处理客户端会阻塞的代码。
func HandleConnection(conn net.Conn) {
	conn.Write([]byte("Welcome connect!"))

	// 读取用户数据
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n') //将r的内容也就是conn的数据按照换行符进行读取。
		if err == io.EOF {
			conn.Close()
		}
		// 除去换行符
		line = strings.TrimSpace(line)
		// 在没有终止信号发出时，等待
		if len(line) == 0 {
			continue
		}

		command := getCommand(line)
		if command["cmd"] == "go get" {
			cmd := cmdcreator.Goget(command["user"], command["package"])
			docker.RunDocker(cmd)
		}
	}
}
