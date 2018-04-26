package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/sysu-go-online/docker_end/cmdcreator"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"

	"github.com/gorilla/websocket"
)

// 异步读取信息，并发送给connection
func writeToConnection(container *Container, hjconn types.HijackedResponse, ctl chan<- bool) {
	for {
		p, _, err := hjconn.Reader.ReadLine()
		container.conn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			ctl <- true
			return
		}
	}
}

// 异步读取信息，并发送给cmd
func readFromClient(dConn net.Conn, cConn *websocket.Conn) {
	defer dConn.Close()

	// Read message from client and write to process
	for {
		_, msg, err := cConn.ReadMessage()
		// If client close connection, kill the process
		if err != nil {
			return
		}
		fmt.Println("Message get: " + string(msg))
		msg = append(msg, byte('\n'))
		_, err = dConn.Write(msg)
		// If message can not be written to the process, kill it
		if err != nil {
			panic(err)
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
		Image: "golang",
		// TODO
		Volumes:    map[string]struct{}{},
		WorkingDir: getPWD(cont.ProjectName, cont.UserName, cont.PWD),
		// Entrypoint: []string{"sh"},
	}
	hostConfig = &container.HostConfig{
		// TODO
		Binds:      []string{},
		AutoRemove: true,
		DNS:        []string{"8.8.8.8"},
	}
	networkingConfig = &network.NetworkingConfig{}
	attachOptions = types.ContainerAttachOptions{
		Stream: true,
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Logs:   true,
	}
	startOptions = types.ContainerStartOptions{}
	return
}

func getPWD(projectname string, username string, pwd string) string {
	goPath := "/go"
	path := filepath.Join(goPath, "src/githb.com/", username, projectname, pwd)
	return path
}
