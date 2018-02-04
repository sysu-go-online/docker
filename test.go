package test

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func mainc() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// _, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "helloworld",
		// AttachStdin: true,
		Tty: true,
		Cmd: []string{"bash"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	ret, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{Stderr: true, Stdin: true, Stdout: true})
	fmt.Println(1)
	msg := []byte{}
	for {
		b := make([]byte, 10)
		if num, _ := ret.Conn.Read(b); num <= 0 {
			break
		}
		msg = append(msg, b...)
	}
	fmt.Println(string(msg))
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
