package cmdcreator

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestCommand_Gocmds(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    *exec.Cmd
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.command.Gocmds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.Gocmds() = %v, want %v", got, tt.want)
			}
		})
	}
}
