package cmd

import (
	"fmt"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	testNsConfig = clientcmdapi.Config{
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
)

func Test_changeNamespace(t *testing.T) {
	type args struct {
		args           []string
		namespaceList  []Namespaces
		currentContext string
		config         *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"ns",
			args{args: []string{"test"},
				namespaceList: []Namespaces{
					{"test", false},
					{"hammer-ns", true}},
				currentContext: "root-context",
				config:         &testNsConfig},
			false,
		},
		{
			"ns-not-exit",
			args{args: []string{"a"},
				namespaceList: []Namespaces{
					{"test", false},
					{"hammer-ns", true}},
				currentContext: "root-context",
				config:         &testNsConfig},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := changeNamespace(tt.args.args, tt.args.namespaceList, tt.args.currentContext, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("changeNamespace() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Printf("Catch ERROR: %v\n", err)
			}
		})
	}
}
