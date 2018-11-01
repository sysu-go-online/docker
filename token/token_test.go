package token

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		userID   string
		token    *string
		callback func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.args.userID, tt.args.token, tt.args.callback)
		})
	}
}

func TestDel(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Del(tt.args.token)
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name  string
		args  args
		want  *User
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetUser(tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
