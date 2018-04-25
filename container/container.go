package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"fmt"

	"github.com/docker/docker/pkg/stdcopy"

	"github.com/docker/docker/client"
	"github.com/sysu-go-online/docker_end/cmdcreator"

	"github.com/gorilla/websocket"
)

const (
	// BufferSize 缓冲区大小，以字节为单位
	BufferSize = 1024
)

var idset int
var DockerClient *client.Client

// Container 通过接口封装输入输出给
type Container struct {
	ID      string              // container ID
	conn    *websocket.Conn     // 绑定的websocket，其中一端
	command *cmdcreator.Command // User command and other messages
}

func init() {
	// Get docker client with preset env
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	DockerClient = dockerClient
}

// NewContainer 新创建一个容器指针
func NewContainer(conn *websocket.Conn, command *cmdcreator.Command) *Container {
	container := Container{
		conn:    conn,
		command: command,
	}
	// Create container
	ctx, config, hostConfig, netwrokingConfig, _, _ := getConfig()
	ret, err := DockerClient.ContainerCreate(ctx, config, hostConfig, netwrokingConfig, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(ret.Warnings)
	ID := ret.ID
	container.ID = ID
	return &container
}

func StartContainer(container *Container) {
	// Attach container
	ctx, _, _, _, attachOptions, startOptions := getConfig()
	hjconn, err := DockerClient.ContainerAttach(ctx, container.ID, attachOptions)
	defer hjconn.Close()
	if err != nil {
		panic(err)
	}
	// Read message from client and send it to docker
	go readFromClient(hjconn.Conn, container.conn)
	// Read message from docker and send it to client
	stdout := WsWriter{
		conn: container.conn,
	}
	stderr := WsWriter{
		conn: container.conn,
	}
	_, err = stdcopy.StdCopy(stdout, stderr, hjconn.Reader)
	if err != nil {
		panic(err)
	}
	// Start container
	err = DockerClient.ContainerStart(ctx, container.ID, startOptions)
	if err != nil {
		panic(err)
	}
}
