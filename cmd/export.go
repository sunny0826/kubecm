package cmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AddCommand add command struct
type ExportCommand struct {
	BaseCommand
}

// Init AddCommand
func (ec *ExportCommand) Init() {
	ec.command = &cobra.Command{
		Use:     "export",
		Short:   "Export the specified context from the kubeconfig",
		Long:    "Export the specified context from the kubeconfig",
		Aliases: []string{"e"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ec.runExport(args)
		},
		Example: exportExample(),
	}
	ec.command.Flags().StringP("file", "f", "", "Path to export kubeconfig files")
	_ = ec.command.MarkFlagRequired("file")
	ec.AddCommands(&DocsCommand{})
}

func (ec *ExportCommand) runExport(args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		confirm, kubeName, err := selectExportContext(config)
		if err != nil {
			return err
		}
		if confirm == "True" {
			config, err = exportContext([]string{kubeName}, config)
			if err != nil {
				return err
			}
		} else {
			return errors.New("nothing exported！")
		}
	} else {
		config, err = exportContext(args, config)
		if err != nil {
			return err
		}
	}

	file, _ := ec.command.Flags().GetString("file")
	err = clientcmd.WriteToFile(*config, file)
	if err != nil {
		return err
	}
	return nil
}

func exportContext(ctxs []string, config *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	var notFinds []string
	exportConfig := clientcmdapi.NewConfig()
	for _, ctx := range ctxs {
		if ec, ok := config.Contexts[ctx]; ok {
			exportConfig.AuthInfos[ec.AuthInfo] = config.AuthInfos[ec.AuthInfo]
			exportConfig.Clusters[ec.Cluster] = config.Clusters[ec.Cluster]
			exportConfig.Contexts[ctx] = config.Contexts[ctx]
			exportConfig.CurrentContext = ctx
			fmt.Printf("Context Export:「%s」\n", ctx)
		} else {
			notFinds = append(notFinds, ctx)
			fmt.Printf("「%s」do not exit.\n", ctx)
		}
	}
	if len(notFinds) == len(ctxs) {
		return nil, errors.New("nothing exported！")
	}
	return exportConfig, nil
}

func selectExportContext(config *clientcmdapi.Config) (string, string, error) {
	var kubeItems []Needle
	for key, obj := range config.Contexts {
		if key != config.CurrentContext {
			kubeItems = append(kubeItems, Needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]Needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	slices.SortFunc(kubeItems, compareKubeItems)
	// exit option
	kubeItems, err := ExitOption(kubeItems)
	if err != nil {
		return "", "", err
	}
	num := SelectUI(kubeItems, "Select The Export Kube Context")
	kubeName := kubeItems[num].Name
	confirm := BoolUI(fmt.Sprintf("Are you sure you want to export「%s」?", kubeName))
	return confirm, kubeName, nil
}

func exportExample() string {
	return `
# Export context to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1
# Export multiple contexts to myconfig.yaml file
kubecm export -f myconfig.yaml my-context1 my-context2
`
}
