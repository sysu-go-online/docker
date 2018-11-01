package container

import (
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/docker_end/cmdcreator"
)

func TestNewContainer(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		command *cmdcreator.Command
	}
	tests := []struct {
		name string
		args args
		want *Container
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContainer(tt.args.conn, tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectNetwork(t *testing.T) {
	type args struct {
		cont *Container
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConnectNetwork(tt.args.cont); (err != nil) != tt.wantErr {
				t.Errorf("ConnectNetwork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartContainer(t *testing.T) {
	type args struct {
		container *Container
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartContainer(tt.args.container)
		})
	}
}

func TestStopContainer(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StopContainer(tt.args.id)
		})
	}
}

func TestUserConf_SetDefault(t *testing.T) {
	type args struct {
		container *Container
	}
	tests := []struct {
		name string
		c    *UserConf
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetDefault(tt.args.container)
		})
	}
}
