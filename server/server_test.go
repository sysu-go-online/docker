package server

import (
	"reflect"
	"testing"

	"github.com/codegangsta/martini"
	"github.com/unrolled/render"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *martini.ClassicMartini
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initRoutes(t *testing.T) {
	type args struct {
		m         *martini.ClassicMartini
		formatter *render.Render
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initRoutes(tt.args.m, tt.args.formatter)
		})
	}
}
