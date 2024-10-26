package cmd

import (
	"errors"
	"fmt"
	"slices"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// SwitchCommand switch cmd struct
type SwitchCommand struct {
	BaseCommand
}

// Init SwitchCommand
func (sc *SwitchCommand) Init() {
	sc.command = &cobra.Command{
		Use:   "switch",
		Short: "Switch Kube Context interactively",
		Long: `
Switch Kube Context interactively
`,
		Aliases: []string{"s", "sw"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("no support for more than 1 parameter")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.runSwitch(cmd, args)
		},
		Example: switchExample(),
	}
	sc.AddCommands(&DocsCommand{})
}

func (sc *SwitchCommand) runSwitch(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	switch len(args) {
	case 0:
		config, err = handleOperation(config)
		if err != nil {
			return err
		}
	case 1:
		config, err = handleQuickSwitch(config, args[0])
		if err != nil {
			return err
		}
	}
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	fmt.Printf("Switched to context 「%s」\n", config.CurrentContext)
	return MacNotifier(fmt.Sprintf("Switched to context [%s]\n", config.CurrentContext))
}

func handleQuickSwitch(config *clientcmdapi.Config, name string) (*clientcmdapi.Config, error) {
	if _, ok := config.Contexts[name]; !ok {
		return config, errors.New("cannot find context named 「" + name + "」")
	}
	config.CurrentContext = name
	return config, nil
}

func handleOperation(config *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	var kubeItems []Needle
	current := config.CurrentContext
	for key, obj := range config.Contexts {
		if key != current {
			kubeItems = append(kubeItems, Needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]Needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	slices.SortFunc(kubeItems, compareKubeItems)
	// exit option
	kubeItems, err := ExitOption(kubeItems)
	if err != nil {
		return config, err
	}
	num := SelectUI(kubeItems, "Select Kube Context")
	kubeName := kubeItems[num].Name
	config.CurrentContext = kubeName
	return config, nil
}

//TODO need update docs

func switchExample() string {
	return `
# Switch Kube Context interactively
kubecm switch
# Quick switch Kube Context
kubecm switch dev
`
}
