package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/sysu-go-online/docker_end/cmdcreator"

	minetypes "github.com/sysu-go-online/docker_end/types"

	"github.com/gorilla/websocket"
)

const (
	// BufferSize 缓冲区大小，以字节为单位
	BufferSize = 1024
)

var idset int
var DockerClient *client.Client
var DefaultLanguage string
var tty bool

// Container 通过接口封装输入输出给
type Container struct {
	ID      string              // container ID
	conn    *websocket.Conn     // 绑定的websocket，其中一端
	command *cmdcreator.Command // User command and other messages
	context *UserConf
}

// UserConf stores conf read from user file
type UserConf struct {
	Language    string
	Username    string
	ProjectName string
	Environment []string
}

func init() {
	// Get docker client with preset env
	dockerClient, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	if err != nil {
		panic(err)
	}
	DockerClient = dockerClient
}

// NewContainer 新创建一个容器指针
// prepare container environment
// read user information from command
// set in-container environment from user-defined file
func NewContainer(msg *minetypes.CreateContainerRequest) string {
	tty = true
	ctx := context.Background()
	config := &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		Env:          []string{},
		Cmd:          []string{"zsh"},
		Image:        msg.Image,
		WorkingDir:   msg.PWD,
	}
	config.Env = append(config.Env, msg.ENV...)
	hostConfig := &container.HostConfig{
		Binds:       []string{},
		AutoRemove:  true,
		DNS:         []string{"8.8.8.8"},
		CapAdd:      []string{"SYS_PTRACE"},
		SecurityOpt: []string{"seccomp=unconfined"},
	}
	for i := range msg.MNT {
		tmp := mount.Mount{
			Type:   mount.TypeBind,
			Source: msg.MNT[i],
			Target: msg.Target[i],
		}
		hostConfig.Mounts = append(hostConfig.Mounts, tmp)
	}
	networkingConfig := &network.NetworkingConfig{}

	ret, err := DockerClient.ContainerCreate(ctx, config, hostConfig, networkingConfig, "")
	if err != nil {
		log.Println(err)
		return ""
	}
	return ret.ID
}

// StartContainer start container
func StartContainer(id string) error {
	return DockerClient.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

// ResizeContainer resize container
func ResizeContainer(id string, width int, height int) error {
	err := StartContainer(id)
	if err != nil {
		return err
	}
	return DockerClient.ContainerResize(context.Background(), id, types.ResizeOptions{
		Height: uint(height),
		Width:  uint(width),
	})
}

// GetHijackRes get attach response
func GetHijackRes(id string) (*types.HijackedResponse, error) {
	r, err := DockerClient.ContainerAttach(context.Background(), id, types.ContainerAttachOptions{
		Stream: true,
		Stderr: true,
		Stdin:  true,
		Stdout: true,
	})
	return &r, err
}

// WriteToUserConn attach to the container
func WriteToUserConn(conn *websocket.Conn, reader *bufio.Reader, judge *bool) {
	for {
		if *judge {
			return
		}
		data := make([]byte, 6)
		n, err := reader.Read(data)
		if err != nil {
			log.Println(err)
			conn.Close()
			*judge = true
			return
		}
		if n == 0 {
			continue
		}
		data = data[:n]
		res := minetypes.ConnectContainerResponse{true, string(data)}
		err = conn.WriteJSON(res)
		if err != nil {
			log.Println(err)
			conn.Close()
			*judge = true
			return
		}
	}
}

// WriteToContainer write data to container
func WriteToContainer(wsconn *websocket.Conn, conconn net.Conn, judge *bool) {
	for {
		if *judge {
			return
		}
		msg := minetypes.ConnectContainerRequest{}
		err := wsconn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			wsconn.Close()
			*judge = true
			return
		}
		_, err = conconn.Write([]byte(msg.Msg))
		if err != nil {
			log.Println(err)
			wsconn.Close()
			*judge = true
			return
		}
	}
}

// ConnectNetwork connect a container to network
func ConnectNetwork(cont *Container) error {
	list, err := DockerClient.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}
	CONTAINERNETWORKNAME := os.Getenv("CONTAINER_NETWORK")
	if len(CONTAINERNETWORKNAME) == 0 {
		CONTAINERNETWORKNAME = "user-services"
	}
	var CONTAINERNETWORKID string
	for _, v := range list {
		if strings.Contains(v.Name, CONTAINERNETWORKNAME) {
			CONTAINERNETWORKID = v.ID
			break
		}
	}
	if CONTAINERNETWORKID == "" {
		fmt.Printf("Can not find network named %s\n", CONTAINERNETWORKNAME)
		return nil
	}
	return DockerClient.NetworkConnect(context.Background(), CONTAINERNETWORKID, cont.ID, &network.EndpointSettings{})
}

// StopContainer stop container after the connection stops
func StopContainer(id string) {
	err := DockerClient.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Println(err)
	}
}

// SetDefault set default value
// TODO: add env according to the language
func (c *UserConf) SetDefault(container *Container) {
	if container != nil && c != nil {
		if c.Language == "" {
			c.Language = DefaultLanguage
		}
		if c.ProjectName == "" {
			c.ProjectName = container.command.ProjectName
		}
		if c.Username == "" {
			c.Username = container.command.UserName
		}
		// add necassary env
		for _, v := range c.Environment {
			firstEqualPos := strings.Index(v, "=")
			if firstEqualPos == -1 {
				continue
			}
			key := v[0:firstEqualPos]
			value := v[firstEqualPos:]
			isSet := false
			for i, v1 := range container.command.ENV {
				firstEqualPos1 := strings.Index(v1, "=")
				if firstEqualPos1 == -1 {
					// remove invalid entry
					container.command.ENV = append(container.command.ENV[:i], container.command.ENV[i+1:]...)
				}
				key1 := v[0:firstEqualPos]
				value1 := v[firstEqualPos:]
				if key1 == key {
					value = strings.Join([]string{value1, value}, ":")
					container.command.ENV[i] = strings.Join([]string{key, value}, "=")
					isSet = true
					break
				}
			}
			if !isSet {
				container.command.ENV = append(container.command.ENV, v)
			}
		}
	}
}
