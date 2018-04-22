package cmdcreator

import "os/exec"

//********************************************
// Author : huziang
//   包含正常cmd的命令实现
//********************************************

// Ls : ls -l
func (command *Command) Ls() *exec.Cmd {
	return exec.Command("docker", "run", "--rm", "-i", "ubuntu", "ls", "-l")
}
