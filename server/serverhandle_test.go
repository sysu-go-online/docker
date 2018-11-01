package server

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/unrolled/render"
)

func TestHandleTTYConnection(t *testing.T) {
	type args struct {
		formatter *render.Render
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleTTYConnection(tt.args.formatter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleTTYConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleDebugConnection(t *testing.T) {
	type args struct {
		formatter *render.Render
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleDebugConnection(tt.args.formatter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleDebugConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
