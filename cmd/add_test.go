package cmd

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	addRootConfigConflictAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"black-user": {Token: "black-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"pig-cluster": {Server: "http://pig.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"root-context": {AuthInfo: "black-user", Cluster: "pig-cluster", Namespace: "saw-ns"}},
	}
	addConfigAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"red-user": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cow-cluster": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"federal-context": {AuthInfo: "red-user", Cluster: "cow-cluster", Namespace: "hammer-ns"}},
	}
	addWantConfigAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-gmbtgkhfch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-gmbtgkhfch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"name": {AuthInfo: "user-gmbtgkhfch", Cluster: "cluster-gmbtgkhfch", Namespace: "hammer-ns"}},
	}
	addTestWantConfigAlfa = clientcmdapi.Config{
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user-gmbtgkhfch": {Token: "red-token"}},
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster-gmbtgkhfch": {Server: "http://cow.org:8080"}},
		Contexts: map[string]*clientcmdapi.Context{
			"test": {AuthInfo: "user-gmbtgkhfch", Cluster: "cluster-gmbtgkhfch", Namespace: "hammer-ns"}},
	}
)

func Test_formatNewConfig(t *testing.T) {
	rootConfig, _ := ioutil.TempFile("", "")
	defer os.Remove(rootConfig.Name())
	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())
	_ = clientcmd.WriteToFile(addRootConfigConflictAlfa, rootConfig.Name())
	_ = clientcmd.WriteToFile(addConfigAlfa, configFile.Name())
	wantName := splitTempName(configFile.Name())
	wantConfig := clientcmdapi.NewConfig()
	addWantConfigAlfa.DeepCopyInto(wantConfig)
	for key, obj := range wantConfig.Contexts {
		wantConfig.Contexts[wantName] = obj
		delete(wantConfig.Contexts, key)
		break
	}
	cfgFile = rootConfig.Name()

	type args struct {
		file     string
		nameFlag string
	}
	tests := []struct {
		name    string
		args    args
		want    *clientcmdapi.Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"name-is-null", args{configFile.Name(), ""}, wantConfig, false},
		{"name-is-set", args{configFile.Name(), "test"}, &addTestWantConfigAlfa, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNewConfig(tt.args.file, tt.args.nameFlag)
			if err != nil {
				return
			}
			checkResult(tt.want, got, t)
		})
	}
}

func Test_formatAndCheckName(t *testing.T) {
	rootConfig, _ := ioutil.TempFile("", "")
	defer os.Remove(rootConfig.Name())
	configFile, _ := ioutil.TempFile("", "")
	defer os.Remove(configFile.Name())
	_ = clientcmd.WriteToFile(addRootConfigConflictAlfa, rootConfig.Name())
	_ = clientcmd.WriteToFile(addConfigAlfa, configFile.Name())
	wantName := splitTempName(configFile.Name())
	cfgFile = rootConfig.Name()

	type args struct {
		file string
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"name-is-null", args{configFile.Name(), ""}, wantName, false},
		{"name-is-set", args{configFile.Name(), "test"}, "test", false},
		{"name-is-exists", args{configFile.Name(), "root-context"}, "root-context", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatAndCheckName(tt.args.file, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatAndCheckName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatAndCheckName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func splitTempName(path string) string {
	name := strings.Split(path, "/")
	wantName := name[len(name)-1]
	return wantName
}

func checkResult(want, got *clientcmdapi.Config, t *testing.T) {
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
