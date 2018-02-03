package cmdcreator

import (
	"os/exec"
	"strings"
)

var GOROOT_ENV = "/usr/lib/go-1.8:/usr/local/go"
var GOPATH_ENV = "/home/huziang/Desktop/home/{0}:/go"

func Goget(username string, packagename string) *exec.Cmd {
	gopath := strings.Replace(GOPATH_ENV, "{0}", username, -1)

	return exec.Command("docker", "run", "--rm", "-i",
		"-v", GOROOT_ENV,
		"-v", gopath,
		"golang", "go", "get", packagename)
}

func Ls() *exec.Cmd {
	return exec.Command("docker", "run", "--rm", "-i", "ubuntu", "ls", "-l")
}
