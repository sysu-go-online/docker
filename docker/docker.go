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
			fmt.Println("?")
			fmt.Println("OUT:" + string(b))
			if err != nil {
				fmt.Println(err)
				outend = true
				close(outchan)
			} else {
				fmt.Println("OUT:" + string(b))
				outchan <- b[:n]
			}
		}

		if !errend {
			n, err := stderr.Read(b)
			fmt.Println("OUT:" + string(b))
			if err != nil {
				errend = true
				close(errchan)
			} else {
				fmt.Println("ERR:" + string(b))
				fmt.Println(n)
				errchan <- b[:n]
			}
		}

		if outend && errend {
			break
		}
	}
	cmd.Wait()
}
