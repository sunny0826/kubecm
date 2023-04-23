package cmd

import (
	"os"
	"testing"
)

func Test_updateFile(t *testing.T) {
	temp, err := os.CreateTemp("", "")
	if err != nil {
		return
	}
	defer os.Remove(temp.Name())
	type args struct {
		cxt  string
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				cxt:  "test",
				path: temp.Name(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateFile(tt.args.cxt, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("updateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.ReadFile(tt.args.path)
			if err != nil {
				t.Errorf("updateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(file) != tt.args.cxt {
				t.Errorf("updateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeAppend(t *testing.T) {
	bashFile, _ := os.CreateTemp("", ".bash")
	defer os.Remove(bashFile.Name())
	type args struct {
		context string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				context: "test",
				path:    bashFile.Name(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeAppend(tt.args.context, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("writeAppend() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.ReadFile(tt.args.path)
			if err != nil {
				t.Errorf("writeAppend() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(file) != tt.args.context+"\n" {
				t.Errorf("writeAppend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
