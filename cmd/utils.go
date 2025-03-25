package cmd

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	r "runtime"
	"sort"
	"strings"
	"time"

	"k8s.io/client-go/rest"

	"github.com/bndr/gotabulate"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/imdario/mergo"
	"github.com/manifoldco/promptui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
)

// Needle use for switch
type Needle struct {
	Name    string
	Cluster string
	User    string
	Center  string
}

// Namespaces namespaces struct
type Namespaces struct {
	Name    string
	Default bool
}

const (
	Filename  = "filename"
	Context   = "context"
	User      = "user"
	Cluster   = "cluster"
	Namespace = "namespace"
)

// SelectRunner interface - For better unit testing
type SelectRunner interface {
	Run() (int, string, error)
}

// PromptRunner interface - For better unit testing
type PromptRunner interface {
	Run() (string, error)
}

// Copied from https://github.com/kubernetes/kubernetes
// /blob/master/pkg/kubectl/util/hash/hash.go
func hEncode(hex string) (string, error) {
	if len(hex) < 10 {
		return "", fmt.Errorf(
			"input length must be at least 10")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}

// Hash returns the hex form of the sha256 of the argument.
func Hash(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

// HashSuf return the string of kubeconfig.
func HashSuf(config *clientcmdapi.Config) string {
	reJSON, err := runtime.Encode(clientcmdlatest.Codec, config)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
	sum, _ := hEncode(Hash(string(reJSON)))
	return sum
}

// HashSufString return the string of HashSuf.
func HashSufString(data string) string {
	sum, _ := hEncode(Hash(data))
	return sum
}

// PrintTable generate table
func PrintTable(config *clientcmdapi.Config) error {
	var table [][]string
	sortedKeys := make([]string, 0)
	for k := range config.Contexts {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	ctx := config.Contexts
	for _, k := range sortedKeys {
		namespace := "default"
		head := ""
		if config.CurrentContext == k {
			head = "*"
		}
		if ctx[k].Namespace != "" {
			namespace = ctx[k].Namespace
		}
		if config.Clusters == nil {
			continue
		}
		cluster, ok := config.Clusters[ctx[k].Cluster]
		if !ok {
			continue
		}
		conTmp := []string{head, k, ctx[k].Cluster, ctx[k].AuthInfo, cluster.Server, namespace}
		table = append(table, conTmp)
	}

	if table != nil {
		tabulate := gotabulate.Create(table)
		tabulate.SetHeaders([]string{"CURRENT", "NAME", "CLUSTER", "USER", "SERVER", "Namespace"})
		// Turn On String Wrapping
		tabulate.SetWrapStrings(true)
		// Render the table
		tabulate.SetAlign("center")
		fmt.Println(tabulate.Render("grid", "left"))
	} else {
		return errors.New("context not found")
	}
	return nil
}

// SelectUI output select ui
func SelectUI(kubeItems []Needle, label string) int {
	s, err := selectUIRunner(kubeItems, label, nil)
	if err != nil {
		if err.Error() == "exit" {
			os.Exit(0)
		}
		log.Fatalf("Prompt failed %v\n", err)
	}
	return s
}

// selectUIRunner
func selectUIRunner(kubeItems []Needle, label string, runner SelectRunner) (int, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F63C {{ .Name | red }}{{ .Center | red}}",
		Inactive: "  {{ .Name | cyan }}{{ .Center | red}}",
		Selected: "\U0001F638 Select:{{ .Name | green }}",
		Details: `
--------- Info ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Cluster:" | faint }}	{{ .Cluster }}
{{ "User:" | faint }}	{{ .User }}`,
	}
	searcher := func(input string, index int) bool {
		pepper := kubeItems[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		if input == "q" && name == "<exit>" {
			return true
		}
		return fuzzy.Match(input, name)
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     kubeItems,
		Templates: templates,
		Size:      uiSize,
		Searcher:  searcher,
	}
	if runner == nil {
		runner = &prompt
	}
	i, _, err := runner.Run()
	if err != nil {
		return 0, err
	}
	if kubeItems[i].Name == "<Exit>" {
		return 0, errors.New("exit")
	}
	return i, err
}

// PromptUI output prompt ui
func PromptUI(label string, name string) string {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("context name must have more than 1 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  name,
	}
	result, err := promptUIWithRunner(&prompt)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// promptUIWithRunner
func promptUIWithRunner(runner PromptRunner) (string, error) {
	result, err := runner.Run()

	if err != nil {
		return "", errors.New("prompt failed")
	}
	return result, nil
}

// BoolUI output bool ui
func BoolUI(label string) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F37A {{ . | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F47B {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     []string{"False", "True"},
		Templates: templates,
		Size:      uiSize,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}

// ClusterStatusCheck check cluster status
type ClusterStatusCheck struct {
	Version   *v.Info
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
}

// ClusterStatus output cluster status
func ClusterStatus(duration time.Duration) (*ClusterStatusCheck, error) {
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		return nil, err
	}
	config.Timeout = time.Second * duration
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	serverVersion, err := clientSet.ServerVersion()
	if err != nil {
		return nil, err
	}

	return &ClusterStatusCheck{
		Version:   serverVersion,
		ClientSet: clientSet,
		Config:    config,
	}, nil
}

// MoreInfo output more info
func MoreInfo(clientSet kubernetes.Interface, writer io.Writer) error {
	if os.Getenv("KUBECM_DISABLE_K8S_MORE_INFO") != "" {
		return nil
	}
	timeout := int64(2)
	ctx := context.TODO()
	nodesList, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}
	podsList, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}
	nsList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return err
	}

	kv := make(map[string]int)
	kv["Namespace"] = len(nsList.Items)
	kv["Node"] = len(nodesList.Items)
	kv["Pod"] = len(podsList.Items)
	printKV(writer, "[Summary] ", kv)
	return nil
}

// WriteConfig write kubeconfig
func WriteConfig(cover bool, file string, outConfig *clientcmdapi.Config) error {
	if cover {
		err := clientcmd.WriteToFile(*outConfig, cfgFile)
		if err != nil {
			return err
		}
		fmt.Printf("「%s」 write successful!\n", file)

		if !silenceTable {
			err = PrintTable(outConfig)
			if err != nil {
				return err
			}
		}

	} else {
		err := clientcmd.WriteToFile(*outConfig, "kubecm.config")
		if err != nil {
			return err
		}
		printString(os.Stdout, "generate ./kubecm.config\n")
	}
	return nil
}

// UpdateConfigFile update kubeconfig
func UpdateConfigFile(file string, updateConfig *clientcmdapi.Config) error {
	file, err := CheckAndTransformFilePath(file, cfgCreate)
	if err != nil {
		return err
	}
	err = clientcmd.WriteToFile(*updateConfig, file)
	if err != nil {
		return err
	}
	printString(os.Stdout, "Update Config: "+file+"\n")
	return nil
}

// ExitOption exit option of SelectUI
func ExitOption(kubeItems []Needle) ([]Needle, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	kubeItems = append(kubeItems, Needle{Name: "<Exit>", Cluster: "exit the kubecm", User: u.Username})
	return kubeItems, nil
}

func printService(out io.Writer, name, link string) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Fprint(out, name)
	ct.ResetColor()
	fmt.Fprint(out, " is running at ")
	ct.ChangeColor(ct.Yellow, false, ct.None, false)
	fmt.Fprint(out, link)
	ct.ResetColor()
	fmt.Fprintln(out, "")
}

func printString(out io.Writer, name string) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Fprint(out, name)
	ct.ResetColor()
}

func printKV(out io.Writer, prefix string, kv map[string]int) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Fprint(out, prefix)
	ct.ResetColor()
	for k, v := range kv {
		ct.ChangeColor(ct.Cyan, false, ct.None, false)
		fmt.Fprint(out, k)
		fmt.Fprint(out, ": ")
		ct.ResetColor()
		ct.ChangeColor(ct.Yellow, false, ct.None, false)
		fmt.Fprint(out, v)
		ct.ResetColor()
		fmt.Fprint(out, " ")
	}
	fmt.Fprint(out, "\n")
}

func printYellow(out io.Writer, content string) {
	ct.ChangeColor(ct.Yellow, false, ct.None, false)
	fmt.Fprint(out, content)
	ct.ResetColor()
}

func printWarning(out io.Writer, name string) {
	ct.ChangeColor(ct.Red, false, ct.None, false)
	fmt.Fprint(out, name)
	ct.ResetColor()
}

func appendConfig(c1, c2 *clientcmdapi.Config) *clientcmdapi.Config {
	config := clientcmdapi.NewConfig()
	_ = mergo.Merge(config, c1)
	_ = mergo.Merge(config, c2)
	return config
}

// CheckAndTransformFilePath return converted path
func CheckAndTransformFilePath(path string, autoCreate bool) (string, error) {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir(), path[2:])
	}
	if IsFile(path) {
		printYellow(os.Stdout, path+" Path Exist\n")
	} else {
		if !autoCreate {
			return path, errors.New("path Not Exist")
		}
		printYellow(os.Stdout, "Createing Directory: "+filepath.Dir(path)+"\n")
		printYellow(os.Stdout, path+" Path is Not Absolute, setting path to home dir\n")
		pathDir := filepath.Join(homeDir(), ".kube")
		path = filepath.Join(pathDir, "config")
		err := os.MkdirAll(pathDir, 0777)
		if err != nil {
			return path, err
		}
		file, err := os.Create(path)
		if err != nil {
			return path, err
		}
		defer file.Close()
		return path, err
	}
	// read files info
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func compareKubeItems(a, b Needle) int {
	return strings.Compare(a.Name, b.Name)
}

// CheckValidContext check and clean mismatched AuthInfo and Cluster
func CheckValidContext(clear bool, config *clientcmdapi.Config) *clientcmdapi.Config {
	for key, obj := range config.Contexts {
		if _, ok := config.AuthInfos[obj.AuthInfo]; !ok {
			if clear {
				printString(os.Stdout, fmt.Sprintf("clear lapsed AuthInfo [%s]\n", obj.AuthInfo))
			} else {
				printYellow(os.Stdout, fmt.Sprintf("WARNING: AuthInfo 「%s」 has no matching context 「%s」, please run `kubecm clear` to clean up this Context.\n", obj.AuthInfo, key))
			}
			delete(config.Contexts, key)
			delete(config.Clusters, obj.Cluster)
		}
		if _, ok := config.Clusters[obj.Cluster]; !ok {
			if clear {
				printString(os.Stdout, fmt.Sprintf("clear lapsed Cluster [%s]\n", obj.Cluster))
			} else {
				printYellow(os.Stdout, fmt.Sprintf("WARNING: Cluster 「%s」 has no matching context 「%s」, please run `kubecm clear` to clean up this Context.\n", obj.Cluster, key))
			}
			delete(config.Contexts, key)
			delete(config.AuthInfos, obj.AuthInfo)
		}
	}
	return config
}

func getFileName(path string) string {
	n := strings.Split(path, "/")
	result := strings.Split(n[len(n)-1], ".")
	return result[0]
}

// MacNotifier send notify message in macOS
func MacNotifier(msg string) error {
	if isMacOs() && macNotify {
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "Kubecm"`, msg))
		return cmd.Run()
	}
	return nil
}

// isMacOs check if current system is macOS
func isMacOs() bool {
	return r.GOOS == "darwin"
}

func validateContextTemplate(contextTemplate []string) error {
	for _, value := range contextTemplate {
		if value != Filename && value != Context && value != User && value != Cluster && value != Namespace {
			return errors.New("the available values for context-template are: filename, user, cluster, context, namespace")
		}
	}
	return nil
}

// checkes if a path exists
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// CheckAndTransformFilePath return converted path
func CheckAndTransformDirPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir(), path[2:])
	}
	// read files info
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
