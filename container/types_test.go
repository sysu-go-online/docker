package container

import "testing"

func TestWsWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		w       WsWriter
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.w.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("WsWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("WsWriter.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
