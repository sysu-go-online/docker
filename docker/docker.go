package docker

import (
	"bufio"
	"io"
	"os/exec"
)

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
