package docker

import (
	"bufio"
	"fmt"
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
	reader := bufio.NewReader(stdout)
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
			line, err := reader.ReadString('\n')
			fmt.Println(line)
			if err != nil {
				outend = true
				close(outchan)
			} else {
				outchan <- []byte(line)
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
