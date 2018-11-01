package util

import "testing"

func TestDealPanic(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DealPanic(tt.args.err)
		})
	}
}

func Test_isFileExit(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFileExit(); got != tt.want {
				t.Errorf("isFileExit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGOPATH(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGOPATH(); got != tt.want {
				t.Errorf("GetGOPATH() = %v, want %v", got, tt.want)
			}
		})
	}
}
