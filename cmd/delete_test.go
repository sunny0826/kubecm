package cmd

import (
	"fmt"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	delMergeConfig = clientcmdapi.Config{
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
	delRootConfigConflictAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"}},
	}
	delButKeepUserConfigBefore = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "black-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
	delButKeepUserConfigAfter = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
		},
	}
)

func Test_deleteContext(t *testing.T) {
	type args struct {
		ctxs   []string
		config *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"delete", args{[]string{"federal-context"}, &delMergeConfig}, false},
		{"delete-not-exist", args{[]string{"a"}, &delMergeConfig}, true},
		{"multiple-delete", args{[]string{"federal-context", "root-context"}, &delMergeConfig}, false},
		{"delete-but-keep-user", args{[]string{"federal-context"}, &delButKeepUserConfigBefore}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteContext(tt.args.ctxs, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("deleteContext() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				switch tt.name {
				case "delete":
					checkConfig(&delRootConfigConflictAlfa, tt.args.config, t)
				case "multiple-delete":
					checkConfig(clientcmdapi.NewConfig(), tt.args.config, t)
				case "delete-but-keep-user":
					checkConfig(&delButKeepUserConfigAfter, tt.args.config, t)
				}
				if err != nil {
					fmt.Println(err)
				}
			}
		})
	}
}
