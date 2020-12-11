package cmd

import (
	"os/user"
	"reflect"
	"testing"

	apiequality "k8s.io/apimachinery/pkg/api/equality"

	"k8s.io/apimachinery/pkg/util/diff"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	appendRootConfigConflictAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"}},
	}
	appendConfigAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"red-user": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cow-cluster": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"}},
	}
	appendMergeConfig = clientcmdapi.Config{
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
	wrongRootConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"red-user": {Token: "red-token"},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cow-cluster": {Server: "http://cow.org:8080"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context":    {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"},
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"},
		},
	}
	wrongFederalConfig = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"},
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

func Test_appendConfig(t *testing.T) {
	type args struct {
		c1 *clientcmdapi.Config
		c2 *clientcmdapi.Config
	}
	tests := []struct {
		name string
		args args
		want *clientcmdapi.Config
	}{
		// TODO: Add test cases.
		{"merge", args{&appendRootConfigConflictAlfa, &appendConfigAlfa}, &appendMergeConfig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendConfig(tt.args.c1, tt.args.c2)
			checkResult(tt.want, got, "", t)
		})
	}
}

func checkResult(want, got *clientcmdapi.Config, wantname string, t *testing.T) {
	testSetNilMapsToEmpties(reflect.ValueOf(&got))
	testSetNilMapsToEmpties(reflect.ValueOf(&want))
	testClearLocationOfOrigin(got)

	if !apiequality.Semantic.DeepEqual(want, got) {
		t.Errorf("diff: %v", diff.ObjectDiff(want, got))
		t.Errorf("expected: %#v\n actual:   %#v", want, got)
	}
}

func testClearLocationOfOrigin(config *clientcmdapi.Config) {
	for key, obj := range config.AuthInfos {
		obj.LocationOfOrigin = ""
		config.AuthInfos[key] = obj
	}
	for key, obj := range config.Clusters {
		obj.LocationOfOrigin = ""
		config.Clusters[key] = obj
	}
	for key, obj := range config.Contexts {
		obj.LocationOfOrigin = ""
		config.Contexts[key] = obj
	}
}

func testSetNilMapsToEmpties(curr reflect.Value) {
	actualCurrValue := curr
	if curr.Kind() == reflect.Ptr {
		actualCurrValue = curr.Elem()
	}

	switch actualCurrValue.Kind() {
	case reflect.Map:
		for _, mapKey := range actualCurrValue.MapKeys() {
			currMapValue := actualCurrValue.MapIndex(mapKey)
			testSetNilMapsToEmpties(currMapValue)
		}

	case reflect.Struct:
		for fieldIndex := 0; fieldIndex < actualCurrValue.NumField(); fieldIndex++ {
			currFieldValue := actualCurrValue.Field(fieldIndex)

			if currFieldValue.Kind() == reflect.Map && currFieldValue.IsNil() {
				newValue := reflect.MakeMap(currFieldValue.Type())
				currFieldValue.Set(newValue)
			} else {
				testSetNilMapsToEmpties(currFieldValue.Addr())
			}
		}

	}

}

func TestExitOption(t *testing.T) {
	gotNeedles := []Needle{
		{"test1", "test2", "any", "*"},
		{"test", "test2", "any", ""},
	}
	u, _ := user.Current()
	wantNeedles := []Needle{
		{"test1", "test2", "any", "*"},
		{"test", "test2", "any", ""},
		{"<Exit>", "exit the kubecm", u.Username, ""},
	}
	type args struct {
		kubeItems []Needle
	}
	tests := []struct {
		name    string
		args    args
		want    []Needle
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", args{gotNeedles}, wantNeedles, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExitOption(tt.args.kubeItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExitOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExitOption() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAndTransformFilePath(t *testing.T) {
	wantPath := homeDir()
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test-~", args{path: "~/"}, wantPath, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAndTransformFilePath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndTransformFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAndTransformFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckValidContext(t *testing.T) {
	type args struct {
		config *clientcmdapi.Config
	}
	tests := []struct {
		name string
		args args
		want *clientcmdapi.Config
	}{
		// TODO: Add test cases.
		{"clear-root", args{config: &wrongRootConfig}, &appendConfigAlfa},
		{"clear-federal", args{config: &wrongFederalConfig}, &appendRootConfigConflictAlfa},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckValidContext(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckValidContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
