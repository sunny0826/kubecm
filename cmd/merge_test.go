package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_listFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("TempDir %s: %v", t.Name(), err)
	}
	defer os.RemoveAll(tempDir)
	filename1 := filepath.Join(tempDir, "config1")
	filename2 := filepath.Join(tempDir, "config2")
	err = ioutil.WriteFile(filename1, []byte("shmorp"), 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", filename1, err)
	}
	err = ioutil.WriteFile(filename2, []byte("florp"), 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", filename2, err)
	}
	type args struct {
		folder string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{"testDir", args{folder: tempDir}, []string{filename1, filename2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listFile(tt.args.folder); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
