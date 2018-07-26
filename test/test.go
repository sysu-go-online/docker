package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// Command 解析客户端发送的信息
type Command struct {
	Command     string   `json:"command"`
	PWD         string   `json:"pwd"`
	ENV         []string `json:"env"`
	UserName    string   `json:"user"`
	ProjectName string   `json:"project"`
	Entrypoint  []string `json:"entrypoint"`
}

// DialDockerService create connection between web server and docker server
func dialDockerService() (*websocket.Conn, error) {
	// Set up websocket connection
	dockerAddr := os.Getenv("DOCKER_ADDRESS")
	dockerPort := os.Getenv("DOCKER_PORT")
	if len(dockerAddr) == 0 {
		dockerAddr = "localhost"
	}
	if len(dockerPort) == 0 {
		dockerPort = "8888"
	}
	dockerPort = ":" + dockerPort
	dockerAddr = dockerAddr + dockerPort
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// HandleMessage decide different operation according to the given json message
func handleMessage(msg string, conn *websocket.Conn, isFirst bool) error {
	var workSpace *Command
	var err error
	if isFirst {
		projectName := "test"
		username := "golang"
		pwd := ""
		env := []string{"GOPATH", "/root/go:/home/go"}
		workSpace = &Command{
			Command:     msg,
			PWD:         pwd,
			ENV:         env,
			UserName:    username,
			ProjectName: projectName,
		}
	}

	// Send message
	if isFirst {
		err = conn.WriteJSON(*workSpace)
	} else {
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
	if err != nil {
		return err
	}
	return nil
}

func main() {
	conn, _ := dialDockerService()
	go func() {
		for {
			t, bs, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			if t == websocket.TextMessage {
				fmt.Print(string(bs))
			}
		}
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			bs, _, _ := reader.ReadLine()
			err := conn.WriteMessage(websocket.TextMessage, bs)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}
	}()
	handleMessage("go run main.go", conn, true)

	for {

	}
}
