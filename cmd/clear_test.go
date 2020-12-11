package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func Test_clearContext(t *testing.T) {
	trueFile, _ := ioutil.TempFile("", "")
	falseFile, _ := ioutil.TempFile("", "")
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
		{"true", args{file: trueFile.Name()}, true, false},
		{"false", args{file: falseFile.Name()}, false, false},
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
