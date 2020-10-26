package cmd

import (
	"fmt"
	"reflect"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	renameWantConfigAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-gmbtgkhfch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-gmbtgkhfch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"test": {AuthInfo: "user-gmbtgkhfch", Cluster: "cluster-gmbtgkhfch", Namespace: "hammer-ns"}},
	}
	renameWantConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-gmbtgkhfch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-gmbtgkhfch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"name": {AuthInfo: "user-gmbtgkhfch", Cluster: "cluster-gmbtgkhfch", Namespace: "hammer-ns"}},
	}
	renameMergeConfig = clientcmdapi.Config{
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
	renameWantConfigB = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-gmbtgkhfch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-gmbtgkhfch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"name": {AuthInfo: "user-gmbtgkhfch", Cluster: "cluster-gmbtgkhfch", Namespace: "hammer-ns"}},
	}
)

func Test_renameComplet(t *testing.T) {
	type args struct {
		rename   string
		kubeName string
		config   *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"rename", args{"test", "name", &renameWantConfigB}, &renameWantConfigAlfa, false},
		{"rename=kubeName", args{"test", "test", &renameWantConfig}, nil, true},
		{"rename-in-config", args{"federal-context", "root-context", &renameMergeConfig}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renameComplet(tt.args.rename, tt.args.kubeName, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("renameComplet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renameComplet() got = %v, want %v", got, tt.want)
			}
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		})
	}
}
