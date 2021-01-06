package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"strconv"
)

// AddCommand add command struct
type AddCommand struct {
	BaseCommand
}

type KubeConfig struct {
	config *clientcmdapi.Config
}

// Init AddCommand
func (ac *AddCommand) Init() {
	ac.command = &cobra.Command{
		Use:   "add",
		Short: "Add KubeConfig to $HOME/.kube/config",
		Long:  "Add KubeConfig to $HOME/.kube/config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.runAdd(cmd, args)
		},
		Example: addExample(),
	}
	ac.command.Flags().StringP("file", "f", "", "Path to merge kubeconfig files")
	_ = ac.command.MarkFlagRequired("file")
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	file, _ := ac.command.Flags().GetString("file")
	// check path
	file, err := CheckAndTransformFilePath(file)
	if err != nil {
		return err
	}
	newConfig, err := clientcmd.LoadFromFile(file)
	if err != nil {
		return err
	}
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	kubeConfig := &KubeConfig{
		config: newConfig,
	}
	// merge context loop
	outConfig, err := kubeConfig.handleContexts(oldConfig)
	if err != nil {
		return err
	}
	if len(outConfig.Contexts) == 1 {
		for k := range outConfig.Contexts {
			outConfig.CurrentContext = k
		}
	}
	cover := BoolUI(fmt.Sprintf("Does it overwrite File 「%s」?", cfgFile))
	confirm, err := strconv.ParseBool(cover)
	if err != nil {
		return err
	}
	err = WriteConfig(confirm, file, outConfig)
	if err != nil {
		return err
	}
	return nil
}

func (kc *KubeConfig) handleContexts(oldConfig *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	newConfig := clientcmdapi.NewConfig()
	for name, ctx := range kc.config.Contexts {
		newName := name
		if checkContextName(name, oldConfig) {
			nameConfirm := BoolUI(fmt.Sprintf("「%s」 Name already exists, do you want to rename it. (If you select `False`, this context will not be merged)", name))
			if nameConfirm == "True" {
				newName = PromptUI("Rename", name)
				if newName == name {
					return nil, errors.New("need to rename")
				}
			} else {
				continue
			}
		}
		itemConfig := kc.handleContext(newName, ctx)
		newConfig = appendConfig(newConfig, itemConfig)
		fmt.Printf("Add Context: %s \n", newName)
	}
	outConfig := appendConfig(oldConfig, newConfig)
	return outConfig, nil
}

func checkContextName(name string, oldConfig *clientcmdapi.Config) bool {
	if _, ok := oldConfig.Contexts[name]; ok {
		return true
	}
	return false
}

func (kc *KubeConfig) handleContext(key string, ctx *clientcmdapi.Context) *clientcmdapi.Config {
	newConfig := clientcmdapi.NewConfig()
	suffix := HashSufString(key)
	userName := fmt.Sprintf("user-%v", suffix)
	clusterName := fmt.Sprintf("cluster-%v", suffix)
	newCtx := ctx.DeepCopy()
	newConfig.AuthInfos[userName] = kc.config.AuthInfos[newCtx.AuthInfo]
	newConfig.Clusters[clusterName] = kc.config.Clusters[newCtx.Cluster]
	newConfig.Contexts[key] = newCtx
	newConfig.Contexts[key].AuthInfo = userName
	newConfig.Contexts[key].Cluster = clusterName
	return newConfig
}

func addExample() string {
	return `
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
`
}
