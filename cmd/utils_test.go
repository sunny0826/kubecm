package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	kubecmVersion "github.com/sunny0826/kubecm/version"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

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

func TestPrintTable(t *testing.T) {
	config := &clientcmdapi.Config{
		Contexts: map[string]*clientcmdapi.Context{
			"test-ctx": {
				Cluster:  "test-cluster",
				AuthInfo: "test-user",
			},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			"test-cluster": {
				Server: "https://very-long-domain-name.example.com",
			},
		},
		CurrentContext: "test-ctx",
	}

	t.Run("ShortServer", func(t *testing.T) {
		var buf bytes.Buffer
		err := PrintTable(&buf, config, &PrintOption{ShortServer: true})
		assert.NoError(t, err)
		// https://very-long-domain-name.example.com (43 chars)
		// Should be truncated to 27 chars + "..."
		// "https://very-long-domain-na..."
		assert.Contains(t, buf.String(), "https://very-long-domain-na...")
	})

	t.Run("NoServer", func(t *testing.T) {
		var buf bytes.Buffer
		err := PrintTable(&buf, config, &PrintOption{NoServer: true})
		assert.NoError(t, err)
		assert.NotContains(t, buf.String(), "SERVER")
		assert.NotContains(t, buf.String(), "https://very-long-domain-name.example.com")
	})

	t.Run("Default", func(t *testing.T) {
		var buf bytes.Buffer
		err := PrintTable(&buf, config, nil)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "SERVER")
		// Due to table wrapping, the full URL might be split. Check for parts.
		assert.Contains(t, buf.String(), "https://very-long-domain-name.")
		assert.Contains(t, buf.String(), "example.com")
	})
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

func TestKubeconfigSplitter(t *testing.T) {
	type args struct {
		kubeconfig string
	}
	tests := []struct {
		name string
		args args
		GoOs string
		want []string
	}{
		{
			name: "one kubeconfig linux",
			GoOs: kubecmVersion.Linux,
			args: args{kubeconfig: "$HOME/Downloads/kubeconfig-6.yaml"},
			want: []string{"$HOME/Downloads/kubeconfig-6.yaml"},
		},
		{
			name: "two kubeconfig linux",
			GoOs: kubecmVersion.Linux,
			args: args{kubeconfig: "/Users/user123/.kube/config:$HOME/Downloads/kubeconfig-6.yaml"},
			want: []string{"/Users/user123/.kube/config", "$HOME/Downloads/kubeconfig-6.yaml"},
		},
		{
			name: "one kubeconfig windows",
			GoOs: kubecmVersion.Windows,
			args: args{kubeconfig: "$HOME/Downloads/kubeconfig-6.yaml"},
			want: []string{"$HOME/Downloads/kubeconfig-6.yaml"},
		},
		{
			name: "two kubeconfig windows",
			GoOs: kubecmVersion.Windows,
			args: args{kubeconfig: "/Users/user123/.kube/config;$HOME/Downloads/kubeconfig-6.yaml"},
			want: []string{"/Users/user123/.kube/config", "$HOME/Downloads/kubeconfig-6.yaml"},
		},
	}
	for _, tt := range tests {
		kubecmVersion.GoOs = tt.GoOs
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, KubeconfigSplitter(tt.args.kubeconfig), "KubeconfigSplitter(%v)", tt.args.kubeconfig)
		})
	}
}

func TestCheckAndTransformFilePath(t *testing.T) {
	wantPath := filepath.Join(homeDir(), ".kube", "config")
	type args struct {
		path      string
		cfgCreate bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test -~ with no auto create", args{path: "~/", cfgCreate: false}, homeDir(), true},
		{"test -~ with auto create enabled", args{path: "~/", cfgCreate: true}, wantPath, false},
		{"test - false config path no auto create", args{path: "", cfgCreate: false}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAndTransformFilePath(tt.args.path, tt.args.cfgCreate)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndTransformFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAndTransformFilePath() got = %v, want %v", got, tt.want)
			}
		})
		if tt.args.cfgCreate {
			t.Cleanup(func() {
				// Remove the file from wantPath after the test run is done
				os.RemoveAll(filepath.Join(homeDir(), ".kube"))
			})
		}
	}
}
func TestCheckAndTransformDirPath(t *testing.T) {
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
		{"test -~ home dir - should pass", args{path: "~/"}, homeDir(), false},
		{"test -~ with auto create enabled", args{path: ""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAndTransformDirPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndTransformDirPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAndTransformDirPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test -~ not a file", args{path: "."}, false},
		{"test - is a file", args{path: "./test.file"}, true},
		{"test - is a config file", args{path: "./config"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				// Create a file at the path
				_ = os.WriteFile(tt.args.path, []byte{}, 0666)
			}
			got := IsFile(tt.args.path)
			if got != tt.want {
				t.Errorf("IsFile() got = %v, want %v", got, tt.want)
			}
		})
		if tt.want {
			t.Cleanup(func() {
				// Remove the file from wantPath after the test run is done
				os.Remove(tt.args.path)
			})
		}
	}
}

func TestCheckValidContext(t *testing.T) {
	clearWrongConfig := wrongFederalConfig.DeepCopy()
	clearWrongWant := appendRootConfigConflictAlfa.DeepCopy()
	type args struct {
		clear  bool
		config *clientcmdapi.Config
	}
	tests := []struct {
		name string
		args args
		want *clientcmdapi.Config
	}{
		// TODO: Add test cases.
		{"check-root", args{clear: false, config: &wrongRootConfig}, &appendConfigAlfa},
		{"check-federal", args{clear: false, config: &wrongFederalConfig}, &appendRootConfigConflictAlfa},
		{"clear-federal", args{clear: true, config: clearWrongConfig}, clearWrongWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckValidContext(tt.args.clear, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckValidContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func checkConfig(want, got *clientcmdapi.Config, t *testing.T) {
	if !apiequality.Semantic.DeepEqual(want, got) {
		t.Errorf("diff: %v", diff.ObjectDiff(want, got))
		t.Errorf("expected: %#v\n actual:   %#v", want, got)
	}
}

func Test_getFileName(t *testing.T) {
	temp, _ := os.CreateTemp("", "kubecm-get-file-")
	defer os.RemoveAll(temp.Name())
	tempFilePath := fmt.Sprintf("%s/%s", temp.Name(), "testPath")
	_ = os.WriteFile(tempFilePath, []byte{}, 0666)

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"TestFileName", args{path: tempFilePath}, "testPath"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileName(tt.args.path); got != tt.want {
				t.Errorf("getFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoreInfo(t *testing.T) {
	// Create a fake client with mock objects
	var clientSet kubernetes.Interface = fake.NewSimpleClientset(
		&corev1.Node{},
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"}},
	)

	// Create a buffer to capture output from printKV
	buf := bytes.Buffer{}

	// Test the MoreInfo function with the fake client
	err := MoreInfo(clientSet, &buf)
	if err != nil {
		t.Errorf("MoreInfo returned an error: %v", err)
	}

	// Check if the output contains expected values
	if strings.Contains(buf.String(), "Namespace: 1") && strings.Contains(buf.String(), "Node: 1") && strings.Contains(buf.String(), "Pod: 1") {
		t.Logf("MoreInfo output is correct")
	} else {
		t.Errorf("MoreInfo output is incorrect: %s", buf.String())
	}
}

type testSelectPrompt struct {
	index int
	err   error
}

func (t *testSelectPrompt) Run() (int, string, error) {
	return t.index, "", t.err
}

func TestSelectUI(t *testing.T) {
	kubeItems := []Needle{
		{Name: "Needle1", Cluster: "Cluster1", User: "User1", Center: "Center1"},
		{Name: "Needle2", Cluster: "Cluster2", User: "User2", Center: "Center2"},
		{Name: "<Exit>", Cluster: "", User: "", Center: ""},
	}

	tests := []struct {
		name          string
		selectPrompt  SelectRunner
		expectedIndex int
		expectError   bool
	}{
		{
			name: "Select Needle1",
			selectPrompt: &testSelectPrompt{
				index: 0,
				err:   nil,
			},
			expectedIndex: 0,
			expectError:   false,
		},
		{
			name: "Select <Exit>",
			selectPrompt: &testSelectPrompt{
				index: 2,
				err:   nil,
			},
			expectedIndex: 0,
			expectError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := selectUIRunner(kubeItems, "Select a needle", test.selectPrompt)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedIndex, index)
			}
		})
	}
}

type testStringPrompt struct {
	result string
	err    error
}

func (t *testStringPrompt) Run() (string, error) {
	return t.result, t.err
}

func TestPromptUI(t *testing.T) {
	tests := []struct {
		name      string
		prompt    *testStringPrompt
		label     string
		expected  string
		expectErr bool
	}{
		{
			name: "Valid input",
			prompt: &testStringPrompt{
				result: "TestName",
				err:    nil,
			},
			label:     "Enter name",
			expected:  "TestName",
			expectErr: false,
		},
		{
			name: "Error occurred",
			prompt: &testStringPrompt{
				result: "",
				err:    errors.New("prompt failed"),
			},
			label:     "Enter name",
			expected:  "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str, err := promptUIWithRunner(test.prompt)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, str)
			}
		})
	}
}

func TestValidateContextTemplate(t *testing.T) {
	type args struct {
		contextTemplate []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid context template",
			args: args{
				contextTemplate: []string{Filename, Context, User, Cluster, Namespace},
			},
			wantErr: false,
		},
		{
			name: "invalid context template",
			args: args{
				contextTemplate: []string{"invalid"},
			},
			wantErr: true,
		},
		{
			name: "empty context template",
			args: args{
				contextTemplate: []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateContextTemplate(tt.args.contextTemplate); (err == nil) == tt.wantErr {
				t.Errorf("validateContextTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
