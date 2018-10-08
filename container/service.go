package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/sysu-go-online/docker_end/cmdcreator"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"

	"github.com/gorilla/websocket"

	"github.com/docker/go-connections/nat"
)

var goImportPath = "/root/go"

// 异步读取信息，并发送给connection
func writeToConnection(container *Container, hjconn types.HijackedResponse, ctl chan<- bool) {
	// if !tty {
	// 	w := WsWriter{
	// 		conn: container.conn,
	// 	}
	// 	for {
	// 		written, err := stdcopy.StdCopy(w, w, hjconn.Reader)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			break
	// 		}
	// 		if written == 0 {
	// 			break
	// 		}
	// 	}
	// } else {
	type ret struct {
		Msg  string `json:"msg"`
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	body := ret{}
	body.ID = container.ID
	body.Type = "tty"
	for {
		p, err := hjconn.Reader.ReadByte()
		if err != nil {
			break
		}
		body.Msg = string(p)
		err = container.conn.WriteJSON(body)
		if err != nil {
			break
		}
	}
	// }
	ctl <- true
}

// 异步读取信息，并发送给cmd
func readFromClient(dConn net.Conn, cConn *websocket.Conn, ctl chan<- bool) {
	defer dConn.Close()

	// Read message from client and write to process
	for {
		msg := &cmdcreator.Command{}
		err := cConn.ReadJSON(msg)
		// If client close connection, kill the process
		if err != nil {
			ctl <- true
			return
		}
		// fmt.Print(string(msg))
		_, err = dConn.Write([]byte(msg.Command))
		// If message can not be written to the process, kill it
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// getConfig returns all the need config with given parameters
// TODO: mount according to language
func getConfig(cont *Container, comm *cmdcreator.Command, tty bool) (ctx context.Context, config *container.Config,
	hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig,
	attachOptions types.ContainerAttachOptions, startOptions types.ContainerStartOptions) {
	ctx = context.Background()
	cmd := strings.Split(cont.command.Command, " ")
	image := getImageName(cont)
	config = &container.Config{
		User:         "root",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          tty,
		OpenStdin:    true,
		Env:          cont.context.Environment,
		Cmd:          cmd,
		Image:        image,
		WorkingDir:   getPWD(cont),
		// Entrypoint:   cont.command.Entrypoint,
	}
	if comm != nil {
		ep := nat.PortSet{}
		for _, v := range comm.Ports {
			p, err := nat.NewPort("tcp", strconv.Itoa(v))
			if err != nil {
				fmt.Println(err)
				continue
			}
			ep[p] = struct{}{}
		}
		config.ExposedPorts = ep
	}
	hostConfig = &container.HostConfig{
		Binds:      []string{},
		AutoRemove: true,
		DNS:        []string{"8.8.8.8"},
		Mounts:     getMountList(cont),
	}
	if cont.command.Type == "debug" {
		hostConfig.CapAdd = []string{"SYS_PTRACE"}
		hostConfig.SecurityOpt = []string{"seccomp=unconfined"}
	}
	networkingConfig = &network.NetworkingConfig{}
	attachOptions = types.ContainerAttachOptions{
		Stream: true,
		Stderr: true,
		Stdout: true,
		Stdin:  true,
		Logs:   false,
	}
	startOptions = types.ContainerStartOptions{}
	return
}

// TODO: return according to the language
func getDestination(cont *Container) string {
	username := cont.context.Username
	projectname := cont.context.ProjectName
	var path string
	if cont.command.Type == "tty" {
		path = filepath.Join("/home", "go/src/github.com", username, projectname)
	} else if cont.command.Type == "debug" {
		path = filepath.Join("/home", username, projectname)
	}
	return path
}

func getPWD(cont *Container) string {
	username := cont.context.Username
	projectname := cont.context.ProjectName
	var path string
	if cont.command.Type == "tty" {
		path = filepath.Join("/home", "go/src/github.com/", username, projectname, cont.command.PWD)
	} else if cont.command.Type == "debug" {
		path = filepath.Join("/home", username, projectname)
	}
	return path
}

// TODO: return according to the language
func getHostDir(cont *Container) string {
	username := cont.command.UserName
	projectname := cont.command.ProjectName
	var path string
	if cont.command.Type == "tty" {
		path = filepath.Join("/home", username, "go/src/github.com", projectname)
	} else if cont.command.Type == "debug" {
		path = filepath.Join("/home", username, "cpp", projectname)
	}
	return path
}

// TODO: return image name by language
func getImageName(container *Container) string {
	if container.command.Type == "debug" {
		return "txzdream/go-online-debug_service:dev"
	} else if container.command.Type == "tty" {
		switch container.context.Language {
		case "golang":
			return "golang"
		}
	}
	// return golang as default
	return "golang"
}

func getMountList(container *Container) []mount.Mount {
	var mounts []mount.Mount
	mounts = append(mounts,
		mount.Mount{
			Type: mount.TypeBind,
			// bind current project only
			Source: getHostDir(container),
			Target: getDestination(container),
		},
		// Mount git config file
		mount.Mount{
			Type:   mount.TypeBind,
			Source: filepath.Join("/home", container.command.UserName, "git/gitconfig"),
			Target: filepath.Join("/etc/gitconfig"),
		},
		mount.Mount{
			Type:   mount.TypeBind,
			Source: filepath.Join("/home", container.command.UserName, "git/.gitconfig"),
			Target: filepath.Join("/home", container.command.UserName, ".gitconfig"),
		})
	if container.command.Type == "tty" {
		mounts = append(mounts, mount.Mount{
			// import path
			Type:   mount.TypeBind,
			Source: filepath.Join("/home", container.command.UserName, "go/import"),
			Target: goImportPath,
		})
	}
	return mounts
}
