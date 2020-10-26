package cmd

import (
	"errors"
	"fmt"
	"os"
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
		Short: "Merge configuration file with $HOME/.kube/config",
		Long:  "Merge configuration file with $HOME/.kube/config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.runAdd(cmd, args)
		},
		Example: addExample(),
	}
	ac.command.Flags().StringP("file", "f", "", "Path to merge kubeconfig files")
	ac.command.Flags().StringP("name", "n", "", "The name of contexts. if this field is null,it will be named with file name.")
	ac.command.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	_ = ac.command.MarkFlagRequired("file")
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	file, _ := ac.command.Flags().GetString("file")
	if fileNotExists(file) {
		return errors.New(file + " file does not exist")
	}
	name := ac.command.Flag("name").Value.String()
	newConfig, err := formatNewConfig(file, name)
	if err != nil {
		return err
	}
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	outConfig := appendConfig(oldConfig, newConfig)
	cover, _ := ac.command.Flags().GetBool("cover")
	err = WriteConfig(cover, file, outConfig)
	if err != nil {
		return err
	}
	return nil
}

func formatNewConfig(file, nameFlag string) (*clientcmdapi.Config, error) {
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		return nil, err
	}
	if len(config.AuthInfos) != 1 {
		return nil, errors.New("Only support add 1 context. You can use `merge` cmd")
	}
	name, err := formatAndCheckName(file, nameFlag)
	if err != nil {
		return nil, err
	}
	suffix := HashSuf(config)
	userName := fmt.Sprintf("user-%v", suffix)
	clusterName := fmt.Sprintf("cluster-%v", suffix)
	for key, obj := range config.AuthInfos {
		config.AuthInfos[userName] = obj
		delete(config.AuthInfos, key)
		break
	}
	for key, obj := range config.Clusters {
		config.Clusters[clusterName] = obj
		delete(config.Clusters, key)
		break
	}
	for key, obj := range config.Contexts {
		obj.AuthInfo = userName
		obj.Cluster = clusterName
		config.Contexts[name] = obj
		delete(config.Contexts, key)
		break
	}
	fmt.Printf("Context Add: %s \n", name)
	return config, nil
}

func formatAndCheckName(file, name string) (string, error) {
	if name == "" {
		n := strings.Split(file, "/")
		result := strings.Split(n[len(n)-1], ".")
		name = result[0]
	}
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return "", err
	}
	for key := range config.Contexts {
		if key == name {
			return key, errors.New("The name: " + name + " already exists, please replace it")
		}
	}
	return name, nil
}

func fileNotExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return !os.IsExist(err)
	}
	return false
}

func addExample() string {
	return `
# Merge 1.yaml with $HOME/.kube/config
kubecm add -f 1.yaml 

# Merge 1.yaml and name contexts test with $HOME/.kube/config
kubecm add -f 1.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f 1.yaml -c
`
}
