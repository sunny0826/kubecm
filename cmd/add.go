/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"strings"
)

type AddCommand struct {
	baseCommand
}

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
	ac.command.MarkFlagRequired("file")
}

func (ac *AddCommand) runAdd(cmd *cobra.Command, args []string) error {
	file, _ := ac.command.Flags().GetString("file")
	if fileNotExists(file) {
		return errors.New(file + " file does not exist")
	}
	newConfig, err := ac.formatNewConfig(file)
	if err != nil {
		return err
	}
	oldConfig, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	outConfig := appendConfig(oldConfig, newConfig)
	cover, _ := ac.command.Flags().GetBool("cover")
	err = ac.WriteConfig(cover, file, outConfig)
	if err != nil {
		return err
	}
	return nil
}

func (ac *AddCommand) formatNewConfig(file string) (*clientcmdapi.Config, error) {
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		return nil, err
	}
	if len(config.AuthInfos) != 1 {
		return nil, errors.New("Only support add 1 context. You can use `merge` cmd.\n")
	}
	name, err := ac.fotmatAndCheckName(file)
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
	return config, nil
}

func (ac *AddCommand) fotmatAndCheckName(file string) (string, error) {
	name, _ := ac.command.Flags().GetString("name")
	if name == "" {
		n := strings.Split(file, "/")
		result := strings.Split(n[len(n)-1], ".")
		name = result[0]
	}
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return "", err
	}
	for key, _ := range config.Contexts {
		if key == name {
			return "", errors.New("The name: " + name + " already exists, please replace it.\n")
		}
	}
	return name, nil
}

func fileNotExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return false
		}
		return true
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
