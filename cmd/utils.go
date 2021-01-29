package cmd

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bndr/gotabulate"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/imdario/mergo"
	"github.com/manifoldco/promptui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:     label,
		Items:     kubeItems,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	if kubeItems[i].Name == "<Exit>" {
		fmt.Println("Exited.")
		os.Exit(1)
	}
	return i
}

// PromptUI output prompt ui
func PromptUI(label string, name string) string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Context name must have more than 3 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  name,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
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
		Size:      2,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return obj
}

// ClusterStatus output cluster status
func ClusterStatus() error {
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		return err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	versionInfo, err := clientSet.ServerVersion()
	if err != nil {
		return err
	}

	printString(os.Stdout, "Cluster check succeeded!")
	printString(os.Stdout, "\nKubernetes version ")
	printYellow(os.Stdout, versionInfo.GitVersion)
	printService(os.Stdout, "\nKubernetes master", config.Host)
	err = MoreInfo(clientSet)
	if err != nil {
		fmt.Println("(Error reporting can be ignored and does not affect usage.)")
	}
	return nil
}

func MoreInfo(clientSet *kubernetes.Clientset) error {
	timeout := int64(5)
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
	printKV(os.Stdout, "[Summary] ", kv)
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
		err = PrintTable(outConfig)
		if err != nil {
			return err
		}
	} else {
		err := clientcmd.WriteToFile(*outConfig, "config.yaml")
		if err != nil {
			return err
		}
		printString(os.Stdout, "generate ./config.yaml\n")
	}
	return nil
}

func UpdateConfigFile(file string, updateConfig *clientcmdapi.Config) error {
	file, err := CheckAndTransformFilePath(file)
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

// GetNamespaceList return namespace list
func GetNamespaceList(cont string) ([]Namespaces, error) {
	var nss []Namespaces
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	ctx := context.TODO()
	namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	for _, specItem := range namespaceList.Items {
		switch cont {
		case "":
			if specItem.Name == "default" {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: true})
			} else {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: false})
			}
		default:
			if specItem.Name == cont {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: true})
			} else {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: false})
			}
		}
	}
	return nss, nil
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
		ct.ChangeColor(ct.Blue, false, ct.None, false)
		fmt.Fprint(out, k)
		fmt.Fprint(out, ": ")
		ct.ResetColor()
		ct.ChangeColor(ct.Yellow, false, ct.None, false)
		fmt.Fprint(out, v)
		ct.ResetColor()
		fmt.Fprint(out, " ")
	}
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

func printComponents(out io.Writer, name string, list []string) {
	ct.ChangeColor(ct.Green, false, ct.None, false)
	fmt.Fprint(out, name)
	ct.ResetColor()
	fmt.Fprint(out, ": ")
	ct.ChangeColor(ct.Yellow, false, ct.None, false)
	fmt.Printf("%v \n", list)
	ct.ResetColor()
	fmt.Fprintln(out, "")
}

func appendConfig(c1, c2 *clientcmdapi.Config) *clientcmdapi.Config {
	config := clientcmdapi.NewConfig()
	_ = mergo.Merge(config, c1)
	_ = mergo.Merge(config, c2)
	return config
}

// CheckAndTransformFilePath return converted path
func CheckAndTransformFilePath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir(), path[2:])
	}
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return "", err
	}
	return path, nil
}

// CheckValidContext check and clean mismatched AuthInfo and Cluster
func CheckValidContext(config *clientcmdapi.Config) *clientcmdapi.Config {
	for key, obj := range config.Contexts {
		if _, ok := config.AuthInfos[obj.AuthInfo]; !ok {
			printString(os.Stdout, "Check Config: ")
			fmt.Printf("AuthInfo 「%s」 has no matching context 「%s」, skip\n", obj.AuthInfo, key)
			delete(config.Contexts, key)
			delete(config.Clusters, obj.Cluster)
		}
		if _, ok := config.Clusters[obj.Cluster]; !ok {
			printString(os.Stdout, "Check Config: ")
			fmt.Printf("Cluster 「%s」 has no matching context 「%s」, skip\n", obj.Cluster, key)
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
