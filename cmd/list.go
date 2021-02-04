package cmd

import (
	"fmt"
	"os"
	"strings"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

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
		Use:     "list",
		Short:   "List KubeConfig",
		Long:    "List KubeConfig",
		Aliases: []string{"ls", "l"},
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
	config = CheckValidContext(false, config)
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
		printWarning(os.Stdout, "Cluster check failure!\n")
		return err
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
		return nil, fmt.Errorf("there is no matching context for %v", args)
	}
	return config, nil
}

func listExample() string {
	return `
# List all the contexts in your KubeConfig file
kubecm list
# Aliases
kubecm ls
kubecm l
# Filter out keywords(Multi-keyword support)
kubecm ls kind k3s
`
}
