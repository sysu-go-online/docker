package container

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/gorilla/websocket"
)

func Test_writeToConnection(t *testing.T) {
	type args struct {
		container *Container
		hjconn    types.HijackedResponse
		ctl       chan<- bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeToConnection(tt.args.container, tt.args.hjconn, tt.args.ctl)
		})
	}
}

func Test_readFromClient(t *testing.T) {
	type args struct {
		dConn net.Conn
		cConn *websocket.Conn
		ctl   chan<- bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readFromClient(tt.args.dConn, tt.args.cConn, tt.args.ctl)
		})
	}
}

func Test_getConfig(t *testing.T) {
	type args struct {
		cont *Container
		tty  bool
	}
	tests := []struct {
		name                 string
		args                 args
		wantCtx              context.Context
		wantConfig           *container.Config
		wantHostConfig       *container.HostConfig
		wantNetworkingConfig *network.NetworkingConfig
		wantAttachOptions    types.ContainerAttachOptions
		wantStartOptions     types.ContainerStartOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx, gotConfig, gotHostConfig, gotNetworkingConfig, gotAttachOptions, gotStartOptions := getConfig(tt.args.cont, tt.args.tty)
			if !reflect.DeepEqual(gotCtx, tt.wantCtx) {
				t.Errorf("getConfig() gotCtx = %v, want %v", gotCtx, tt.wantCtx)
			}
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("getConfig() gotConfig = %v, want %v", gotConfig, tt.wantConfig)
			}
			if !reflect.DeepEqual(gotHostConfig, tt.wantHostConfig) {
				t.Errorf("getConfig() gotHostConfig = %v, want %v", gotHostConfig, tt.wantHostConfig)
			}
			if !reflect.DeepEqual(gotNetworkingConfig, tt.wantNetworkingConfig) {
				t.Errorf("getConfig() gotNetworkingConfig = %v, want %v", gotNetworkingConfig, tt.wantNetworkingConfig)
			}
			if !reflect.DeepEqual(gotAttachOptions, tt.wantAttachOptions) {
				t.Errorf("getConfig() gotAttachOptions = %v, want %v", gotAttachOptions, tt.wantAttachOptions)
			}
			if !reflect.DeepEqual(gotStartOptions, tt.wantStartOptions) {
				t.Errorf("getConfig() gotStartOptions = %v, want %v", gotStartOptions, tt.wantStartOptions)
			}
		})
	}
}

func Test_getDestination(t *testing.T) {
	type args struct {
		cont *Container
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDestination(tt.args.cont); got != tt.want {
				t.Errorf("getDestination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPWD(t *testing.T) {
	type args struct {
		cont *Container
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPWD(tt.args.cont); got != tt.want {
				t.Errorf("getPWD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHostDir(t *testing.T) {
	type args struct {
		cont *Container
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostDir(tt.args.cont); got != tt.want {
				t.Errorf("getHostDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getImageName(t *testing.T) {
	type args struct {
		container *Container
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImageName(tt.args.container); got != tt.want {
				t.Errorf("getImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMountList(t *testing.T) {
	type args struct {
		container *Container
	}
	tests := []struct {
		name string
		args args
		want []mount.Mount
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMountList(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMountList() = %v, want %v", got, tt.want)
			}
		})
	}
}
