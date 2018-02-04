package docker

import (
	"bufio"
	"io"
	"os/exec"
)

// RunDocker 运行容器
func RunDocker(cmd *exec.Cmd, in *io.PipeReader, outchan chan []byte, errchan chan []byte) {
	// docker stdin, stdout, stderr设置
	cmd.Stdin = bufio.NewReader(in)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	defer stderr.Close()

	// docker 启动
	cmd.Start()

	// docker write
	b := make([]byte, 1024)
	outend := false
	errend := false
	for {
		if !outend {
			n, err := stdout.Read(b)
			if err != nil {
				outend = true
				close(outchan)
			} else {
				outchan <- b[:n]
			}
		}

		if !errend {
			n, err := stderr.Read(b)
			if err != nil {
				errend = true
				close(errchan)
			} else {
				errchan <- b[:n]
			}
		}

		if outend && errend {
			break
		}
	}
	cmd.Wait()
}
