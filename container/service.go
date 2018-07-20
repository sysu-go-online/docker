package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"

	"github.com/gorilla/websocket"
)

var goImportPath = "/root/go"

// 异步读取信息，并发送给connection
func writeToConnection(container *Container, hjconn types.HijackedResponse, ctl chan<- bool) {
	// w := WsWriter{
	// 	conn: container.conn,
	// }
	// _, err := stdcopy.StdCopy(w, w, hjconn.Reader)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// ctl <- true
	for {
		p, err := hjconn.Reader.ReadByte()
		container.conn.WriteMessage(websocket.TextMessage, []byte{p})
		if err != nil {
			ctl <- true
			return
		}
	}
}

// 异步读取信息，并发送给cmd
func readFromClient(dConn net.Conn, cConn *websocket.Conn, ctl chan<- bool) {
	defer dConn.Close()

	// Read message from client and write to process
	for {
		_, msg, err := cConn.ReadMessage()
		// If client close connection, kill the process
		if err != nil {
			ctl <- true
			return
		}
		// fmt.Println(string(msg))
		// msg = append(msg, byte('\n'))
		_, err = dConn.Write(msg)
		// If message can not be written to the process, kill it
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// getConfig returns all the need config with given parameters
// TODO: mount according to language
func getConfig(cont *Container) (ctx context.Context, config *container.Config,
	hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig,
	attachOptions types.ContainerAttachOptions, startOptions types.ContainerStartOptions) {
	ctx = context.Background()
	cmd := strings.Split(cont.command.Command, " ")
	image := getImageName(cont.context.Language)
	config = &container.Config{
		User:         "root",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		Env:          cont.context.Environment,
		Cmd:          cmd,
		Image:        image,
		WorkingDir:   getPWD(cont.context.ProjectName, cont.context.Username, cont.command.PWD),
		// Entrypoint: []string{"sh"},
	}
	hostConfig = &container.HostConfig{
		Binds:      []string{},
		AutoRemove: true,
		DNS:        []string{"8.8.8.8"},
		Mounts: []mount.Mount{
			// project
			mount.Mount{
				Type: mount.TypeBind,
				// bind current project only
				Source: getHostDir(cont.command.ProjectName, cont.command.UserName),
				Target: getDestination(cont.context.ProjectName, cont.context.Username),
			},
			// import path
			mount.Mount{
				Type:   mount.TypeBind,
				Source: filepath.Join("/home", cont.command.UserName, "go/import"),
				Target: goImportPath,
			},
		},
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
func getDestination(projectname string, username string) string {
	path := filepath.Join("/home", "go/src/github.com", username, projectname)
	fmt.Println(path)
	return path
}

func getPWD(projectname string, username string, pwd string) string {
	path := filepath.Join("/home", "go/src/github.com/", username, projectname, pwd)
	return path
}

// TODO: return according to the language
func getHostDir(projectname string, username string) string {
	home := "/home"
	path := filepath.Join(home, username, "go/src/github.com", projectname)
	return path
}

// TODO: return image name by language
func getImageName(language string) string {
	// Return golang by default
	return "golang"
}
