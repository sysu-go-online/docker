package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"

	"github.com/docker/docker/api/types"
	"github.com/sysu-go-online/docker_end/cmdcreator"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"

	"github.com/gorilla/websocket"
	// . "github.com/sysu-go-online/docker_end/util"
)

var goPath string = "/root/go"
// 异步读取信息，并发送给connection
func writeToConnection(container *Container, hjconn types.HijackedResponse, ctl chan<- bool) {
	// w := WsWriter{
	// 	conn: container.conn,
	// }
	// _, err := stdcopy.StdCopy(w, w, hjconn.Reader)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// ctl <- true
	for {
		p, err := hjconn.Reader.ReadByte()
		container.conn.WriteMessage(websocket.TextMessage, []byte{p})
		if err != nil {
			ctl <- true
			return
		}
	}
}

// 异步读取信息，并发送给cmd
func readFromClient(dConn net.Conn, cConn *websocket.Conn, ctl chan<- bool) {
	defer dConn.Close()

	// Read message from client and write to process
	for {
		_, msg, err := cConn.ReadMessage()
		// If client close connection, kill the process
		if err != nil {
			ctl <- true
			return
		}
		// fmt.Println(string(msg))
		// msg = append(msg, byte('\n'))
		_, err = dConn.Write(msg)
		// If message can not be written to the process, kill it
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// getConfig returns all the need config with given parameters
func getConfig(cont *cmdcreator.Command) (ctx context.Context, config *container.Config,
	hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig,
	attachOptions types.ContainerAttachOptions, startOptions types.ContainerStartOptions) {
	ctx = context.Background()
	cmd := strings.Split(cont.Command, " ")
	config = &container.Config{
		User:         "root",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		Env:          cont.ENV,
		Cmd:          cmd,
		// TODO
		Image:      "golang",
		WorkingDir: getPWD(cont.ProjectName, cont.UserName, cont.PWD),
		// Entrypoint: []string{"sh"},
	}
	hostConfig = &container.HostConfig{
		Binds:      []string{},
		AutoRemove: true,
		DNS:        []string{"8.8.8.8"},
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: getHostDir(cont.ProjectName, cont.UserName),
				Target: getDestination(),
			},
		},
	}
	networkingConfig = &network.NetworkingConfig{}
	attachOptions = types.ContainerAttachOptions{
		Stream: true,
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Logs:   false,
	}
	startOptions = types.ContainerStartOptions{}
	return
}

func getDestination() string {
	return goPath
}

func getPWD(projectname string, username string, pwd string) string {
	path := filepath.Join(goPath, "src/github.com/", username, projectname, pwd)
	return path
}

func getHostDir(projectname string, username string) string {
	// 使用ini文件动态读取配置环境
	// file, err := ini.LoadFile(filepath.Join(GetGOPATH(), "/src/github.com/sysu-go-online/docker_end/config.ini"))
	// if err != nil {
	// 	panic(err)
	// }
	// home, ok := file.Get("HostInformation", "home")
	// if !ok {
	// 	panic(errors.New("读取配置文件发生错误"))
	// }
	home := "/Users/huziang/Desktop/home"
	path := filepath.Join(home, username, "go")
	return path
}
