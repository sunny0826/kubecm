package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	mergeTestConfig = clientcmdapi.Config{
		APIVersion: "v1",
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
)

func Test_listFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("TempDir %s: %v", t.Name(), err)
	}
	defer os.RemoveAll(tempDir)
	filename1 := filepath.Join(tempDir, "config1")
	filename2 := filepath.Join(tempDir, "config2")
	dsStore := filepath.Join(tempDir, ".DS_Store")
	err = ioutil.WriteFile(filename1, []byte("shmorp"), 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", filename1, err)
	}
	err = ioutil.WriteFile(filename2, []byte("florp"), 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", filename2, err)
	}
	err = ioutil.WriteFile(dsStore, []byte("xxxx"), 0444)
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

func Test_loadKubeConfig(t *testing.T) {
	tempDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("TempDir %s: %v", t.Name(), err)
	}
	defer os.RemoveAll(tempDir)

	merge1 := filepath.Join(tempDir, "merge1")
	err = clientcmd.WriteToFile(mergeTestConfig, merge1)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", merge1, err)
	}
	mergeFail := filepath.Join(tempDir, "config2")
	err = ioutil.WriteFile(mergeFail, []byte("florp"), 0444)
	if err != nil {
		t.Fatalf("WriteFile %s: %v", mergeFail, err)
	}

	resultConfig, err := clientcmd.LoadFromFile(merge1)
	if err != nil {
		t.Fatalf("getConfig %s: %v", merge1, err)
	}
	type args struct {
		yaml string
	}
	tests := []struct {
		name    string
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		{
			name: "success config",
			args: args{
				yaml: merge1,
			},
			want:    resultConfig,
			wantErr: false,
		},
		{
			name: "get err file",
			args: args{
				yaml: mergeFail,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadKubeConfig(tt.args.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadKubeConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadKubeConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
