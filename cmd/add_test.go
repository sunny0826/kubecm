package cmd

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	addTestConfig = clientcmdapi.Config{
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
	oldTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
	handleConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-cbc897d6ch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-cbc897d6ch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "user-cbc897d6ch", Cluster: "cluster-cbc897d6ch", Namespace: "hammer-ns"}},
	}
	mergedConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user":      {Token: "black-token"},
			"red-user":        {Token: "red-token"},
			"user-cbc897d6ch": {Token: "red-token"},
			"user-d2m9fd8b7d": {Token: "black-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster":        {Server: "http://pig.org:8080"},
			"cow-cluster":        {Server: "http://cow.org:8080"},
			"cluster-cbc897d6ch": {Server: "http://cow.org:8080"},
			"cluster-d2m9fd8b7d": {Server: "http://pig.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root":            {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal":         {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
			"root-context":    {AuthInfo: "user-d2m9fd8b7d", Cluster: "cluster-d2m9fd8b7d", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "user-cbc897d6ch", Cluster: "cluster-cbc897d6ch", Namespace: "hammer-ns"},
		},
	}
)

func Test_checkContextName(t *testing.T) {
	type args struct {
		name      string
		oldConfig *clientcmdapi.Config
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add more test cases.
		{"exit", args{name: "root-context", oldConfig: &addTestConfig}, true},
		{"not-exit", args{name: "test", oldConfig: &addTestConfig}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkContextName(tt.args.name, tt.args.oldConfig); got != tt.want {
				t.Errorf("checkContextName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeConfig_handleContext(t *testing.T) {
	newConfig := addTestConfig.DeepCopy()
	testCtx := clientcmdapi.Context{AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"}

	type fields struct {
		config *clientcmdapi.Config
	}
	type args struct {
		key string
		ctx *clientcmdapi.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *clientcmdapi.Config
	}{
		// TODO: Add more test cases.
		{"one", fields{config: newConfig}, args{"federal-context", &testCtx}, &handleConfig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KubeConfig{
				config: tt.fields.config,
			}
			got := kc.handleContext(tt.args.key, tt.args.ctx)
			checkConfig(got, tt.want, t)
		})
	}
}

func TestKubeConfig_handleContexts(t *testing.T) {
	newConfig := addTestConfig.DeepCopy()
	type fields struct {
		config *clientcmdapi.Config
	}
	type args struct {
		oldConfig *clientcmdapi.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", fields{config: newConfig}, args{&oldTestConfig}, &mergedConfig, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KubeConfig{
				config: tt.fields.config,
			}
			got, err := kc.handleContexts(tt.args.oldConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleContexts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			checkConfig(got, tt.want, t)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("handleContexts() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
