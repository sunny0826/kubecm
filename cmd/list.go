package cmd

import (
	"errors"
	"fmt"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// ListCommand list cmd struct
type ListCommand struct {
	BaseCommand
}

// Init ListCommand
func (lc *ListCommand) Init() {
	lc.command = &cobra.Command{
		Use:     "ls",
		Short:   "List kubeconfig",
		Long:    "List kubeconfig",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.runList(cmd, args)
		},
		Example: listExample(),
	}
	lc.command.DisableFlagsInUseLine = true
}

func (lc *ListCommand) runList(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	config = CheckValidContext(config)
	outConfig, err := filterArgs(args, config)
	if err != nil {
		return err
	}
	err = PrintTable(outConfig)
	if err != nil {
		return err
	}
	err = ClusterStatus()
	if err != nil {
		return errors.New(fmt.Sprintf("Cluster check failure!\n%v", err))
	}
	return nil
}

func filterArgs(args []string, config *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	if len(args) == 0 {
		return config, nil
	}
	contextList := make(map[string]string)
	for key := range config.Contexts {
		for _, search := range args {
			if strings.Contains(key, search) {
				contextList[key] = search
			}
		}
	}
	for key := range config.Contexts {
		if _, ok := contextList[key]; !ok {
			delete(config.Contexts, key)
		}
	}
	if len(config.Contexts) == 0 {
		return nil, errors.New(fmt.Sprintf("There is no matching context for %v\n", args))
	}
	return config, nil
}

func listExample() string {
	return `
# List all the contexts in your kubeconfig file
kubecm ls
# Aliases
kubecm l
# Filter out keywords(Multi-keyword support)
kubecm ls kind k3s
`
}
