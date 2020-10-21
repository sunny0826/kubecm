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
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"strings"

	"github.com/spf13/cobra"
)

type MergeCommand struct {
	baseCommand
}

func (mc *MergeCommand) Init() {
	mc.command = &cobra.Command{
		Use:     "merge",
		Short:   "Merge the kubeconfig files in the specified directory",
		Long:    `Merge the kubeconfig files in the specified directory`,
		Aliases: []string{"m"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return mc.runMerge(cmd, args)
		},
		Example: mergeExample(),
	}
	mc.command.Flags().StringP("folder", "f", "", "Kubeconfig folder")
	mc.command.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	mc.command.MarkFlagRequired("folder")
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	folder, _ := mc.command.Flags().GetString("folder")
	files := listFile(folder)
	mc.command.Printf("Loading kubeconfig file: %v \n", files)
	configs := clientcmdapi.NewConfig()
	for _, yaml := range files {
		config, err := clientcmd.LoadFromFile(yaml)
		if err != nil {
			return err
		}
		name := nameHandle(yaml)
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
		configs.CurrentContext = name
		configs = appendConfig(configs, config)
		mc.command.Printf("Context Add: %s \n", name)
	}
	cover, _ := mc.command.Flags().GetBool("cover")
	err := mc.WriteConfig(cover, folder, configs)
	if err != nil {
		return err
	}
	return nil
}

func listFile(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var flist []string
	for _, file := range files {
		if file.IsDir() {
			listFile(folder + "/" + file.Name())
		} else {
			flist = append(flist, fmt.Sprintf("%s/%s", folder, file.Name()))
		}
	}
	return flist
}

func nameHandle(path string) string {
	n := strings.Split(path, "/")
	result := strings.Split(n[len(n)-1], ".")
	return result[0]
}

func mergeExample() string {
	return `
# Merge kubeconfig in the dir directory
kubecm merge -f dir

# Merge kubeconfig in the directory and overwrite the original kubeconfig file
kubecm merge -f dir -c
`
}
