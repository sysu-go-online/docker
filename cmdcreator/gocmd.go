package cmdcreator

//********************************************
// Author : huziang
//   包含go语言cmd的实现
//********************************************

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//. "github.com/sysu-go-online/docker_end/util"
)

// GOENV .
const (
	usersHome = "/home"
	gopath    = "/Users/huziang/Desktop/go"
)

// Gocmds : go comannd
func (command *Command) Gocmds() *exec.Cmd {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// test username and project name
	// _, err := os.Stat(filepath.Join(usersHome, command.UserName, "src/github.com", command.ProjectName, command.PWD))
	// DealPanic(err)

	// set mount point
	mountpoints := []string{filepath.Join(usersHome, command.UserName, "src/github.com", command.ProjectName, command.PWD) +
		":" + filepath.Join("/go", "src/github.com", command.ProjectName)}
	fgopath := strings.Split(os.Getenv("GOPATH"), ":")[0]
	mountpoints = append(mountpoints, fgopath+":"+"/go")

	// set work path
	workpath := filepath.Join("/go", "src/github.com", command.ProjectName)

	// set envirment
	envirment := []string{}
	for i := 0; i < len(command.ENV); i += 2 {
		envirment = append(envirment, command.ENV[i]+"="+command.ENV[i+1])
	}

	// set all paramete
	strs := []string{"run", "--rm", "-i"}
	strs = append(strs, "-v")
	strs = append(strs, mountpoints...)
	strs = append(strs, "--env")
	strs = append(strs, envirment...)
	strs = append(strs, "--workdir")
	strs = append(strs, workpath)
	strs = append(strs, "golang")
	strs = append(strs, strings.Split(command.Command, " ")...)
	fmt.Println(strs)
	return exec.Command("docker", strs...)
}
