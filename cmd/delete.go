package cmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// DeleteCommand delete cmd struct
type DeleteCommand struct {
	BaseCommand
}

// Init DeleteCommand
func (dc *DeleteCommand) Init() {
	dc.command = &cobra.Command{
		Use:     "delete",
		Short:   "Delete the specified context from the kubeconfig",
		Long:    `Delete the specified context from the kubeconfig`,
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.runDelete(cmd, args)
		},
		Example: deleteExample(),
	}
	dc.AddCommands(&DocsCommand{})
}

func (dc *DeleteCommand) runDelete(command *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		confirm, kubeName, err := selectDeleteContext(config)
		if err != nil {
			return err
		}
		if confirm == "True" {
			err = deleteContext([]string{kubeName}, config)
			if err != nil {
				return err
			}
		} else {
			return errors.New("nothing deleted！")
		}
	} else {
		err = deleteContext(args, config)
		if err != nil {
			return err
		}
	}
	err = WriteConfig(true, cfgFile, config)
	if err != nil {
		return err
	}
	return nil
}

func deleteContext(ctxs []string, config *clientcmdapi.Config) error {
	var notFinds []string
	for _, ctx := range ctxs {
		if _, ok := config.Contexts[ctx]; ok {
			delContext := config.Contexts[ctx]
			isClusterNameExist, isUserNameExist := checkClusterAndUserNameExceptContextToDelete(config, config.Contexts[ctx])
			if !isUserNameExist {
				delete(config.AuthInfos, delContext.AuthInfo)
			}
			if !isClusterNameExist {
				delete(config.Clusters, delContext.Cluster)
			}
			delete(config.Contexts, ctx)
			fmt.Printf("Context Delete:「%s」\n", ctx)
		} else {
			notFinds = append(notFinds, ctx)
			fmt.Printf("「%s」do not exit.\n", ctx)
		}
	}
	if len(notFinds) == len(ctxs) {
		return errors.New("nothing deleted！")
	}
	return nil
}

func checkClusterAndUserNameExceptContextToDelete(oldConfig *clientcmdapi.Config, contextToDelete *clientcmdapi.Context) (bool, bool) {
	var (
		isClusterNameExist bool
		isUserNameExist    bool
	)

	for _, ctx := range oldConfig.Contexts {
		if ctx.Cluster == contextToDelete.Cluster && ctx != contextToDelete {
			isClusterNameExist = true
		}
		if ctx.AuthInfo == contextToDelete.AuthInfo && ctx != contextToDelete {
			isUserNameExist = true
		}
	}

	return isClusterNameExist, isUserNameExist
}

func selectDeleteContext(config *clientcmdapi.Config) (string, string, error) {
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
	num := SelectUI(kubeItems, "Select The Delete Kube Context")
	kubeName := kubeItems[num].Name
	confirm := BoolUI(fmt.Sprintf("Are you sure you want to delete「%s」?", kubeName))
	return confirm, kubeName, nil
}

func deleteExample() string {
	return `
# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
# Deleting multiple contexts
kubecm delete my-context1 my-context2
`
}
