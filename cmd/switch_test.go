package cmd

import (
	"reflect"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	switchConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"}},
	}
	currentContextSwitchConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"}},
		CurrentContext: "root-context",
	}
)

func Test_handleQuickSwitch(t *testing.T) {
	type args struct {
		config *clientcmdapi.Config
		name   string
	}
	tests := []struct {
		name    string
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"exist-context", args{&switchConfig, "root-context"}, &currentContextSwitchConfig, false},
		{"no-exist-context", args{&switchConfig, "test"}, &currentContextSwitchConfig, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleQuickSwitch(tt.args.config, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleQuickSwitch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleQuickSwitch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
