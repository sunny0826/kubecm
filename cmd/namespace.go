package cmd

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// NamespaceCommand namespace cmd struct
type NamespaceCommand struct {
	BaseCommand
}

// Init NamespaceCommand
func (nc *NamespaceCommand) Init() {
	nc.command = &cobra.Command{
		Use:   "namespace",
		Short: "Switch or change namespace interactively",
		Long: `
Switch or change namespace interactively
`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"ns"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nc.runNamespace(cmd, args)
		},
		Example: namespaceExample(),
	}
	nc.AddCommands(&DocsCommand{})
}

func (nc *NamespaceCommand) runNamespace(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	if len(config.Contexts) == 0 {
		return fmt.Errorf("no valid context found")
	}

	currentContext := config.CurrentContext
	currentNamespace := config.Contexts[currentContext].Namespace
	clientset, err := GetClientSet(cfgFile)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		namespaceList, err := GetNamespaceList(currentNamespace, clientset)
		if err != nil {
			return err
		}
		// exit option
		namespaceList = append(namespaceList, Namespaces{Name: "<Exit>", Default: false})
		num := selectNamespace(namespaceList)
		config.Contexts[currentContext].Namespace = namespaceList[num].Name
	} else {
		exist, err := CheckNamespaceExist(args[0], clientset)
		if err != nil {
			return errors.New("Can not find namespace: " + args[0])
		}
		if exist {
			config.Contexts[currentContext].Namespace = args[0]
			fmt.Printf("Namespace: 「%s」 is selected.\n", args[0])
		} else {
			return errors.New("Can not find namespace: " + args[0])
		}
	}
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	return MacNotifier(fmt.Sprintf("Switch to the [%s] namespace\n", config.Contexts[currentContext].Namespace))
}

func selectNamespace(namespaces []Namespaces) int {
	ns, err := selectNamespaceWithRunner(namespaces, nil)
	if err != nil {
		if err.Error() == "exit" {
			os.Exit(0)
		}
		log.Fatalf("Prompt failed %v\n", err)
	}
	return ns
}

func selectNamespaceWithRunner(namespaces []Namespaces, runner SelectRunner) (int, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F6A9 {{if .Default}} {{ .Name | red }} * {{else}} {{ .Name | red }} {{end}}",
		Inactive: "{{if .Default}} {{ .Name | cyan }} * {{else}} {{ .Name | cyan }} {{end}}",
		Selected: "\U0001F680" + `{{if ne .Name "<Exit>" }}  Namespace: {{ .Name | green }} is selected.{{end}}`,
	}
	searcher := func(input string, index int) bool {
		pepper := namespaces[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		if input == "q" && name == "<exit>" {
			return true
		}
		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:     "Select Namespace:",
		Items:     namespaces,
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
	if namespaces[i].Name == "<Exit>" {
		return 0, errors.New("exit")
	}
	return i, err
}

// GetClientSet return clientset
func GetClientSet(configFile string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return kubernetes.NewForConfig(config)
}

// GetNamespaceList return namespace list
func GetNamespaceList(currentNamespace string, clientset kubernetes.Interface) ([]Namespaces, error) {
	var nss []Namespaces
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	for _, specItem := range namespaceList.Items {
		switch currentNamespace {
		case "":
			if specItem.Name == "default" {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: true})
			} else {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: false})
			}
		default:
			if specItem.Name == currentNamespace {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: true})
			} else {
				nss = append(nss, Namespaces{Name: specItem.Name, Default: false})
			}
		}
	}
	return nss, nil
}

// CheckNamespaceExist check namespace exist
func CheckNamespaceExist(namespace string, clientset kubernetes.Interface) (bool, error) {
	ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf(err.Error())
	}
	if ns.Name == namespace {
		return true, nil
	}
	return false, errors.New("namespace not found")
}

func namespaceExample() string {
	return `
# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system
`
}
