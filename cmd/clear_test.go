package cmd

import (
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func Test_clearContext(t *testing.T) {
	trueFile, _ := os.CreateTemp("", "")
	falseFile, _ := os.CreateTemp("", "")
	defer os.Remove(trueFile.Name())
	defer os.Remove(falseFile.Name())
	_ = clientcmd.WriteToFile(appendMergeConfig, trueFile.Name())
	_ = clientcmd.WriteToFile(wrongRootConfig, falseFile.Name())

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Clear Context", args{file: trueFile.Name()}, true, false},
		//{"Not Clear Context", args{file: falseFile.Name()}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clearContext(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("clearContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("clearContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}
