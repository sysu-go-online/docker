package container

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types/strslice"

	"github.com/gorilla/websocket"
)

// 异步读取信息，并发送给connection
func writeToConnection(out io.ReadCloser, conn *websocket.Conn) {
	defer out.Close()
	reader := bufio.NewReader(out)
	for {
		line, _, err := reader.ReadLine()
		fmt.Println("Message send: " + string(line))
		// 传输结束，发出信号
		if err != nil {
			break
		} else {
			conn.WriteMessage(websocket.TextMessage, line)
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
			panic(err)
			break
		}
		fmt.Println("Message get: " + string(msg))
		msg = append(msg, byte('\n'))
		_, err = dConn.Write(msg)
		// If message can not be written to the process, kill it
		if err != nil {
			panic(err)
			break
		}
	}
}

// getConfig returns all the need config with given parameters
func getConfig() (ctx context.Context, config *container.Config,
	hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig,
	attachOptions types.ContainerAttachOptions, startOptions types.ContainerStartOptions) {
	ctx = context.Background()
	config = &container.Config{
		User:         "root",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		// TODO
		Env: []string{},
		// TODO
		Cmd: strslice.StrSlice{},
		// TODO
		Image: "golang",
		// TODO
		Volumes: map[string]struct{}{},
		// TODO
		WorkingDir: "/",
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
