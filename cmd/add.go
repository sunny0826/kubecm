package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// AddCommand add command struct
type AddCommand struct {
	BaseCommand
}

// Init AddCommand
func (ac *AddCommand) Init() {
	ac.command = &cobra.Command{
		Use:   "add",
		Short: "Add kubeconfig to $HOME/.kube/config",
		Long:  "Add kubeconfig to $HOME/.kube/config",
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
	file, err := CheckAndTransformFilePath(file)
	if err != nil {
		return err
	}
	newConfig, newName, err := formatNewConfig(file)
	if err != nil {
		return err
	}
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	outConfig := appendConfig(oldConfig, newConfig)
	if len(outConfig.Contexts) == 1 {
		for k := range outConfig.Contexts {
			outConfig.CurrentContext = k
		}
	}
	cover := BoolUI(fmt.Sprintf("Are you sure you want to add 「%s」 to the 「%s」context?", newName, cfgFile))
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

func formatNewConfig(file string) (*clientcmdapi.Config, string, error) {
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		return nil, "", err
	}
	//if len(config.AuthInfos) != 1 {
	//	return nil, "", errors.New("Only support add 1 context. You can use `merge` cmd")
	//}
	name, err := formatAndCheckName(file)
	config = CheckValidContext(config)
	if err != nil {
		return nil, "", err
	}
	newConfig := clientcmdapi.NewConfig()
	for key, ctx := range config.Contexts {
		c := handleContext(key, ctx, config)
		appendConfig(newConfig,c)
		fmt.Printf("Context Add: %s \n", key)
	}
	//suffix := HashSuf(config)
	//userName := fmt.Sprintf("user-%v", suffix)
	//clusterName := fmt.Sprintf("cluster-%v", suffix)
	//for key, obj := range config.AuthInfos {
	//	config.AuthInfos[userName] = obj
	//	delete(config.AuthInfos, key)
	//	break
	//}
	//for key, obj := range config.Clusters {
	//	config.Clusters[clusterName] = obj
	//	delete(config.Clusters, key)
	//	break
	//}
	//for key, obj := range config.Contexts {
	//	obj.AuthInfo = userName
	//	obj.Cluster = clusterName
	//	config.Contexts[name] = obj
	//	delete(config.Contexts, key)
	//	break
	//}
	//fmt.Printf("Context Add: %s \n", name)
	return newConfig, name, nil
}

//TODO 支持多 context，将格式化逻辑拆出

func handleContext(key string, ctx *clientcmdapi.Context, config *clientcmdapi.Config) *clientcmdapi.Config {
	newConfig := clientcmdapi.NewConfig()
	suffix := HashSufString(key)
	userName := fmt.Sprintf("user-%v", suffix)
	clusterName := fmt.Sprintf("cluster-%v", suffix)
	newCtx := ctx.DeepCopy()
	newConfig.AuthInfos[userName] = config.AuthInfos[newCtx.AuthInfo]
	newConfig.Clusters[clusterName] = config.Clusters[newCtx.Cluster]
	newConfig.Contexts[key] = newCtx
	newConfig.Contexts[key].AuthInfo = userName
	newConfig.Contexts[key].Cluster = clusterName
	return newConfig
}

func formatAndCheckName(file string) (string, error) {
	n := strings.Split(file, "/")
	result := strings.Split(n[len(n)-1], ".")
	name := result[0]
	nameConfirm := BoolUI(fmt.Sprintf("Need to rename 「%s」 context?", name))
	if nameConfirm == "True" {
		name = PromptUI("Rename", name)
	}
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return "", err
	}
	for key := range config.Contexts {
		if key == name {
			return key, errors.New("The name: 「" + name + "」 already exists, please select another one.")
		}
	}
	return name, nil
}

func addExample() string {
	return `
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
`
}
