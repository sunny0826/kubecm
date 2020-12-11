package cmd

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"reflect"
	"testing"
)

var (
	noRootMergeConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
	noFederalMergeConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
		},
	}
)

func NewAllConfig() *clientcmdapi.Config {
	return appendMergeConfig.DeepCopy()
}

func Test_filterArgs(t *testing.T) {

	type args struct {
		args   []string
		config *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"noFederal", args{[]string{"root"}, NewAllConfig()}, &noFederalMergeConfig, false},
		{"noRoot", args{[]string{"federal"}, NewAllConfig()}, &noRootMergeConfig, false},
		{"all-context", args{[]string{"root", "federal"}, NewAllConfig()}, &appendMergeConfig, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterArgs(tt.args.args, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterArgs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
