package cmd

import (
	"fmt"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	exportTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
	WantExportConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
		},
		CurrentContext: "root-context",
	}
	wantMultipleExportConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
		CurrentContext: "federal-context",
	}
)

func Test_exportContext(t *testing.T) {
	type args struct {
		ctxs   []string
		config *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"export", args{[]string{"root-context"}, &exportTestConfig}, false},
		{"export-not-exist", args{[]string{"a"}, &exportTestConfig}, true},
		{"multiple-export", args{[]string{"root-context", "federal-context"}, &exportTestConfig}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exportedConfig, err := exportContext(tt.args.ctxs, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("exportContext() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				switch tt.name {
				case "export":
					checkConfig(&WantExportConfig, exportedConfig, t)
				case "multiple-export":
					checkConfig(&wantMultipleExportConfig, exportedConfig, t)
				}
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}
