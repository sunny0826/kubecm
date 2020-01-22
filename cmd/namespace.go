package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

type NamespaceCommand struct {
	baseCommand
}

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
}

func (nc *NamespaceCommand) runNamespace(command *cobra.Command, args []string) error {
	config, err := LoadClientConfig(cfgFile)
	if err != nil {
		return nil
	}
	currentContext := config.CurrentContext
	contNs := config.Contexts[currentContext].Namespace
	namespaceList, err := GetNamespaceList(contNs)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		// exit option
		namespaceList = append(namespaceList, namespaces{Name: "<Exit>", Default: false})
		num := selectNamespace(namespaceList)
		config.Contexts[currentContext].Namespace = namespaceList[num].Name
	} else {
		var flag bool
		for _, ns := range namespaceList {
			if ns.Name == args[0] {
				config.Contexts[currentContext].Namespace = args[0]
				nc.command.Printf("Namespace: 「%s」 is selected.\n", args[0])
				flag = true
				break
			}
		}
		if !flag {
			nc.command.Printf("Can not find namespace: 「%s」\n", args[0])
			os.Exit(1)
		}
	}
	err = ModifyKubeConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func selectNamespace(namespaces []namespaces) int {
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
		Size:      4,
		Searcher:  searcher,
	}
	i, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	if namespaces[i].Name == "<Exit>" {
		fmt.Println("Exited.")
		os.Exit(1)
	}
	return i
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
