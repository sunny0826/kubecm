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
	shortServer bool
	noServer    bool
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
	lc.command.Flags().BoolVar(&lc.shortServer, "short-server", false, "Shorten the server endpoint")
	lc.command.Flags().BoolVar(&lc.noServer, "no-server", false, "Hide the server column")
	lc.AddCommands(&DocsCommand{})
}

func (lc *ListCommand) runList(command *cobra.Command, args []string) error {
	clusterMessageChan := make(chan *ClusterStatusCheck)
	go func() {
		info, _ := ClusterStatus(2)
		clusterMessageChan <- info
	}()
	clusterMessage := <-clusterMessageChan
	for _, kubeconfig := range KubeconfigSplitter(cfgFile) {

		config, err := clientcmd.LoadFromFile(kubeconfig)
		if err != nil {
			return err
		}
		config = CheckValidContext(false, config)
		outConfig, err := filterArgs(args, config)
		if err != nil {
			return err
		}
		err = PrintTable(os.Stdout, outConfig, &PrintOption{
			ShortServer: lc.shortServer,
			NoServer:    lc.noServer,
		})
		if err != nil {
			return err
		}

		if clusterMessage != nil {
			printString(os.Stdout, "Cluster check succeeded!")
			printString(os.Stdout, "\nKubernetes version ")
			printYellow(os.Stdout, clusterMessage.Version.GitVersion)
			printService(os.Stdout, "\nKubernetes master", clusterMessage.Config.Host)
			if err := MoreInfo(clusterMessage.ClientSet, os.Stdout); err != nil {
				fmt.Println("(Error reporting can be ignored and does not affect usage.)")
			}
		}
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
# Useful environment variables
KUBECM_DISABLE_K8S_MORE_INFO: it will disable the k8s more info in the output
`
}
