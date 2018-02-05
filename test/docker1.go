package docker1

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main1() {
	fmt.Println("hello?")
	ch := make(chan []byte, 16)

	go func() {
		cmd := exec.Command("docker", "run", "--rm", "-i", "ubuntu", "ls", "-l")
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
