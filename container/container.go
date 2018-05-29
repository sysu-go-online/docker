package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
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

	// Get config
	ctx, config, hostConfig, netwrokingConfig, _, _ := getConfig(container.command)

	// find image
	imagename := "golang"
	images, err := DockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	find := false
	for _, image := range images {
		if strings.Split(image.RepoTags[0], ":")[0] == imagename {
			find = true
			break
		}
	}

	// if not find, pull image
	if !find {
		_, err := DockerClient.ImagePull(ctx, imagename, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
	}

	// Create container
	ret, err := DockerClient.ContainerCreate(ctx, config, hostConfig, netwrokingConfig, "")

	if err != nil {
		panic(err)
	}
	container.ID = ret.ID
	return &container
}

// StartContainer attach and start a container
func StartContainer(container *Container) {
	defer StopContainer(container.ID)
	// Attach container
	ctx, _, _, _, attachOptions, startOptions := getConfig(container.command)
	hjconn, err := DockerClient.ContainerAttach(ctx, container.ID, attachOptions)
	defer hjconn.Close()
	if err != nil {
		panic(err)
	}
	readCtl := make(chan bool, 2)
	// Read message from client and send it to docker
	go readFromClient(hjconn.Conn, container.conn, readCtl)
	// Read message from docker and send it to client
	go writeToConnection(container, hjconn, readCtl)
	// Start container
	err = DockerClient.ContainerStart(ctx, container.ID, startOptions)
	if err != nil {
		fmt.Println(err)
	}
	<-readCtl
}

// StopContainer stop container after the connection stops
func StopContainer(id string) {
	duration := time.Duration(time.Second * 2)
	err := DockerClient.ContainerStop(context.Background(), id, &duration)
	if err != nil {
		fmt.Println(err)
	}
}
