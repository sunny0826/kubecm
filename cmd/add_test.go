package cmd

import (
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
	handleConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-cbc897d6ch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-cbc897d6ch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "user-cbc897d6ch", Cluster: "cluster-cbc897d6ch", Namespace: "hammer-ns"}},
	}
)

//func Test_handleContext(t *testing.T) {
//	wantConfig:= handleConfig.DeepCopy()
//	type args struct {
//		key    string
//		ctx    *clientcmdapi.Context
//		config *clientcmdapi.Config
//	}
//	tests := []struct {
//		name string
//		args args
//		want *clientcmdapi.Config
//	}{
//		// TODO: Add test cases.
//		{"test", args{key: "federal-context", ctx: addTestConfig.Contexts["federal-context"], config: &addTestConfig}, wantConfig},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := handleContext(tt.args.key, tt.args.ctx, tt.args.config)
//			checkConfig(got,tt.want,t)
//		})
//	}
//}
