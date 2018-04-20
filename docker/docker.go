package docker

import (
	"bufio"
	"context"
	"io"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sysu-go-online/docker/cmdcreator"
)

var cli *client.Client
var err error

func init() {
	cli, err = client.NewEnvClient()
}

// 异步读取给定的Reader
func asynWriteChannel(out io.ReadCloser, ch chan []byte) {
	defer out.Close()
	defer close(ch)
	reader := bufio.NewReader(out)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		} else {
			ch <- []byte(line)
		}
	}
}

// RunDocker 运行容器
func RunDocker(cmd *exec.Cmd, in *io.PipeReader, outchan chan []byte, errchan chan []byte) {
	// docker stdin, stdout, stderr设置
	cmd.Stdin = bufio.NewReader(in)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// docker 启动
	cmd.Start()

	// asyn write two channel
	go asynWriteChannel(stdout, outchan)
	go asynWriteChannel(stderr, errchan)
}

// RunCommand receives command and send it to container to run
func RunCommand(commands <-chan cmdcreator.Command, dockerMsg chan<- []byte) {
	isFirst := true
	var containerID string
	var connection types.HijackedResponse
	for command := range commands {
		// if it is the first time to run, init the container first
		if isFirst {
			containerID, connection = initContainer(command)
			go func(dockerMsg chan<- []byte) {
				for msg, err := connection.Reader.ReadBytes('\n'); err == nil; {
					dockerMsg <- msg
				}
			}(dockerMsg)
			cli.ContainerStart(context.Background(), containerID, )
			isFirst = false
		}
	}
}

func initContainer(command cmdcreator.Command) (string, types.HijackedResponse) {
	// Create container
	config := container.Config{
		User:         "root",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		Env:          command.ENV,
		Cmd:          strslice.StrSlice{command.Command},
		Image:        "ubuntu",
		// TODO
		WorkingDir:      "/root",
		NetworkDisabled: false,
	}
	hostConfig := container.HostConfig{
		AutoRemove: true,
	}
	networkConfig := network.NetworkingConfig{}
	ret, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, &networkConfig, "")
	if err != nil {
		panic(err)
	}
	if len(ret.ID) == 0 {
		panic(ret.Warnings)
	}
	// Attach container
	attachConfig := types.ContainerAttachOptions{
		Stream: true,
		Stdin: true,
		Stdout: true,
		Stderr: true,
	}
	conn, err := cli.ContainerAttach(context.Background(), ret.ID, attachConfig)
	if err != nil {
		panic(err)
	}
	return ret.ID, conn
}
