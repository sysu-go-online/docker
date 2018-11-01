package container

import (
	"reflect"
	"testing"
)

func TestNewSignal(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want *Signal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSignal(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSignal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignal_Signal(t *testing.T) {
	tests := []struct {
		name   string
		signal *Signal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.signal.Signal()
		})
	}
}

func TestSignal_Wait(t *testing.T) {
	tests := []struct {
		name   string
		signal *Signal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.signal.Wait()
		})
	}
}
