package cmdcreator

import (
	"os/exec"
	"strings"
)

// GOENV .
const (
	usersHome = "/home/huziang/Desktop/home"
	// goRootEnv = "/usr/lib/go-1.8:/usr/local/go"
	goPathEnv = usersHome + "/{0}:/go"
)

// Goget : go get url
func (command *Command) Goget() *exec.Cmd {
	// gopath := strings.Replace(goPathEnv, "{0}", command.UserName, -1)

	strs := strings.Split(command.Command, " ")
	strs = append([]string{"run", "--rm", "-i", "golang"}, strs...)
	return exec.Command("docker", strs...)
}

// Ls : ls -l
func Ls() *exec.Cmd {
	return exec.Command("docker", "run", "--rm", "-i", "ubuntu", "ls", "-l")
}
