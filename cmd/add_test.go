package cmd

import (
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	addTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
			"red-user":   {Token: "red-token"},
			"not-exist":  {Token: "not-exist-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"},
			"cow-cluster": {Server: "http://cow.org:8080"},
			"not-exist":   {Server: "http://not.exist:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":      {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context":   {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
			"not-exist-context": {AuthInfo: "not-exist", Cluster: "not-exist", Namespace: "not-exist-ns"},
		},
		CurrentContext: "root-context",
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
			"red-user-cbc897d6ch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cow-cluster-cbc897d6ch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "red-user-cbc897d6ch", Cluster: "cow-cluster-cbc897d6ch", Namespace: "hammer-ns"},
		},
	}
	handleNotExistConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"not-exist": {Token: "not-exist-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"not-exist": {Server: "http://not.exist:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"not-exist-context": {AuthInfo: "not-exist", Cluster: "not-exist", Namespace: "not-exist-ns"},
		},
	}
	mergedConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user":            {Token: "black-token"},
			"red-user":              {Token: "red-token"},
			"red-user-7f65b9cc8f":   {Token: "red-token"},
			"black-user-gtch2cf96d": {Token: "black-token"},
			"not-exist":             {Token: "not-exist-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster":            {Server: "http://pig.org:8080"},
			"cow-cluster":            {Server: "http://cow.org:8080"},
			"cow-cluster-7f65b9cc8f": {Server: "http://cow.org:8080"},
			"pig-cluster-gtch2cf96d": {Server: "http://pig.org:8080"},
			"not-exist":              {Server: "http://not.exist:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root":            {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal":         {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
			"test-d2m9fd8b7d": {AuthInfo: "black-user-gtch2cf96d", Cluster: "pig-cluster-gtch2cf96d", Namespace: "saw-ns"},
			"test-cbc897d6ch": {AuthInfo: "red-user-7f65b9cc8f", Cluster: "cow-cluster-7f65b9cc8f", Namespace: "hammer-ns"},
			"test-2h6782585t": {AuthInfo: "not-exist", Cluster: "not-exist", Namespace: "not-exist-ns"},
		},
	}
	singleTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"single-user": {Token: "single-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"single-cluster": {Server: "http://single:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"single-context": {AuthInfo: "single-user", Cluster: "single-cluster", Namespace: "single-ns"},
		},
	}
	mergeSingleTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user":  {Token: "black-token"},
			"red-user":    {Token: "red-token"},
			"single-user": {Token: "single-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster":    {Server: "http://pig.org:8080"},
			"cow-cluster":    {Server: "http://cow.org:8080"},
			"single-cluster": {Server: "http://single:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root":           {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal":        {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
			"single-context": {AuthInfo: "single-user", Cluster: "single-cluster", Namespace: "single-ns"},
		},
	}
	renameSingleTestConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user":  {Token: "black-token"},
			"red-user":    {Token: "red-token"},
			"single-user": {Token: "single-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster":    {Server: "http://pig.org:8080"},
			"cow-cluster":    {Server: "http://cow.org:8080"},
			"single-cluster": {Server: "http://single:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
			"rename":  {AuthInfo: "single-user", Cluster: "single-cluster", Namespace: "single-ns"},
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
		{"exist", args{name: "root-context", oldConfig: &addTestConfig}, true},
		{"not-exist", args{name: "test", oldConfig: &addTestConfig}, false},
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
	testNotExistCtx := clientcmdapi.Context{AuthInfo: "not-exist", Cluster: "not-exist", Namespace: "not-exist-ns"}

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
		{"two", fields{config: newConfig}, args{"not-exist-context", &testNotExistCtx}, &handleNotExistConfig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KubeConfigOption{
				config: tt.fields.config,
			}
			got := kc.handleContext(&oldTestConfig, tt.args.key, tt.args.ctx)
			checkConfig(tt.want, got, t)
		})
	}
}

func TestKubeConfig_handleContexts(t *testing.T) {
	newConfig := addTestConfig.DeepCopy()
	singleConfig := singleTestConfig.DeepCopy()
	type fields struct {
		config   *clientcmdapi.Config
		fileName string
	}
	type args struct {
		oldConfig *clientcmdapi.Config
		newName   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"not have new context name", fields{config: newConfig, fileName: "test"}, args{&oldTestConfig, ""}, &mergedConfig, false},
		{"single context name", fields{config: singleConfig, fileName: "test"}, args{&oldTestConfig, ""}, &mergeSingleTestConfig, false},
		{"single context name - new", fields{config: singleConfig, fileName: "test"}, args{&oldTestConfig, "rename"}, &renameSingleTestConfig, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KubeConfigOption{
				config:   tt.fields.config,
				fileName: tt.fields.fileName,
			}
			got, err := kc.handleContexts(tt.args.oldConfig, tt.args.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleContexts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			checkConfig(tt.want, got, t)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("handleContexts() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestAddToLocal(t *testing.T) {
	localFile, err := os.CreateTemp("", "local-kubeconfig-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	cfgFile = "test"
	// Create a new temporary file
	tempFile, err := os.CreateTemp("", "temp-kubeconfig-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer os.Remove(localFile.Name())

	// Write an initial empty config to the temp file
	emptyConfig := clientcmdapi.NewConfig()
	err = clientcmd.WriteToFile(*emptyConfig, tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to write empty config to temp file: %v", err)
	}
	tempFile.Close()

	err = clientcmd.WriteToFile(addTestConfig, localFile.Name())
	if err != nil {
		t.Fatalf("Failed to write empty config to temp file: %v", err)
	}
	localFile.Close()

	cfgFile = localFile.Name()

	// Mock configuration
	newConfig := &clientcmdapi.Config{
		Clusters:       map[string]*clientcmdapi.Cluster{"test-cluster": {Server: "https://test-cluster"}},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{"test-authinfo": {Token: "black-token"}},
		Contexts:       map[string]*clientcmdapi.Context{"test-context": {AuthInfo: "test-authinfo", Cluster: "test-cluster", Namespace: "hammer-ns"}},
		CurrentContext: "test-context",
	}

	// Test AddToLocal function
	err = AddToLocal(newConfig, tempFile.Name(), "", true)
	if err != nil {
		t.Fatalf("Failed to add to local: %v", err)
	}

	// Read the file and check if the new configuration is added
	loadedConfig, err := clientcmd.LoadFromFile(localFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	if _, ok := loadedConfig.Contexts["test-context"]; !ok {
		t.Fatalf("Failed to find 'test-context' in the loaded config")
	}
}
