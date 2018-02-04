package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("123")
	ch := make(chan []byte, 16)

	go func() {
		cmd := exec.Command("docker", "run", "--rm", "-i", "helloworld")
		cmd.Stdin = bufio.NewReader(os.Stdin)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		for {
			b := make([]byte, 100)
			_, err := stdout.Read(b)
			if err != nil {
				break
			}
			ch <- b
		}
		close(ch)
	}()

	for {
		ob, ok := <-ch
		if ok == false {
			break
		}
		fmt.Print(string(ob))
	}
}
