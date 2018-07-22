package container

//********************************************
// Author : huziang
//   封装容器类，将connection和cmd封装到一起，去除中
// 间chan转换过程
//********************************************

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

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
func NewContainer(conn *websocket.Conn, command *cmdcreator.Command) *Container {
	container := Container{
		conn:    conn,
		command: command,
	}

	// **********Read information from user***********
	userProjectConfPath := filepath.Join("/home", command.UserName, command.ProjectName, "go-online.yml")
	userProjectConf := UserConf{}
	if _, err := os.Stat(userProjectConfPath); os.IsExist(err) {
		userProjectConfData, err := ioutil.ReadFile(userProjectConfPath)
		if err != nil {
			fmt.Println(err)
		}
		if err = yaml.Unmarshal(userProjectConfData, &userProjectConf); err != nil {
			fmt.Println(err)
		}
	}
	userProjectConf.SetDefault(&container)
	container.context = &userProjectConf
	// ***********************************************

	// **********get type and decide image************
	var imagename string
	switch command.Type {
	case "tty":
		imagename = "golang"
		tty = true
	case "debug":
		imagename = "txzdream/go-online-debug_service:dev"
		tty = false
	}
	// ***********************************************

	// Get config
	ctx, config, hostConfig, netwrokingConfig, _, _ := getConfig(&container, tty)

	// find image
	// TODO: match image tag
	images, err := DockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	find := false
	for _, image := range images {
		if len(image.RepoTags) > 0 && strings.Split(image.RepoTags[0], ":")[0] == strings.Split(imagename, ":")[0] {
			find = true
			break
		}
	}

	// if not find, pull image
	if !find {
		fmt.Println("Cant not find such image, trying to pull it")
		_, err := DockerClient.ImagePull(ctx, imagename, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
	}

	// Create container, if image is not pull, wait 10s once until the image pull
	// ret, err := DockerClient.ContainerCreate(ctx, config, hostConfig, netwrokingConfig, "")
	s3, _ := time.ParseDuration("3s")
	for {
		ret, err := DockerClient.ContainerCreate(ctx, config, hostConfig, netwrokingConfig, "")
		if err != nil && strings.Contains(err.Error(), "No such image") {
			time.Sleep(s3)
			continue
		} else if err != nil {
			panic(err)
		} else {
			container.ID = ret.ID
			break
		}
	}

	return &container
}

// StartContainer attach and start a container
func StartContainer(container *Container) {
	defer StopContainer(container.ID)
	// Attach container
	ctx, _, _, _, attachOptions, startOptions := getConfig(container, false)
	hjconn, err := DockerClient.ContainerAttach(ctx, container.ID, attachOptions)
	defer hjconn.Close()
	if err != nil {
		panic(err)
	}
	readCtl := make(chan bool, 2)
	// Read message from client and send it to docker
	go readFromClient(hjconn.Conn, container.conn, readCtl)
	// Read message from docker and send it to client
	go writeToConnection(container, hjconn, readCtl, tty)
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
