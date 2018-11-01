package cmdcreator

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestCommand_Ls(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    *exec.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.Ls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.Ls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    *exec.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.Test(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
