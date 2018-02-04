package main

import (
	"fmt"
	"os/exec"
)

func test(inch chan []byte, outch chan []byte, cmd *exec.Cmd) {

	stdout, _ := cmd.StdoutPipe()
	bs := make([]byte, 1024)
	cmd.Start()
	for {
		n, err := stdout.Read(bs)
		if err != nil {
			close(outch)
			break
		}
		outch <- bs[:n]
	}
	cmd.Wait()
}

func main3() {
	cmd := exec.Command("docker", "run", "-i", "--rm", "ubuntu", "ls")
	outch := make(chan []byte, 1024)
	inch := make(chan []byte, 1024)
	go test(inch, outch, cmd)
	for {
		ob, ok := <-outch
		if ok == false {
			break
		}
		fmt.Print(string(ob))
	}
}
