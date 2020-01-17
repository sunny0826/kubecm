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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	"os"
	"sigs.k8s.io/yaml"
	"strings"

	"github.com/spf13/cobra"
)

type MergeCommand struct {
	baseCommand
}

var folder string

func (mc *MergeCommand) Init() {
	mc.command = &cobra.Command{
		Use:     "merge",
		Short:   "Merge the kubeconfig files in the specified directory",
		Long:    `Merge the kubeconfig files in the specified directory`,
		Example: mergeExample(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return mc.runMerge(cmd, args)
		},
	}
	mc.command.Flags().StringVarP(&folder, "folder", "f", "", "Kubeconfig folder")
	mc.command.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	mc.command.MarkFlagRequired("folder")
}

func (mc MergeCommand) runMerge(command *cobra.Command, args []string) error {
	cover, _ = mc.command.Flags().GetBool("cover")
	files := listFile(folder)
	mc.command.Printf("Loading kubeconfig file: %v \n", files)
	var loop []string
	var defaultName string
	for _, yaml := range files {
		config, err := LoadClientConfig(yaml)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		name := nameHandle(yaml)
		commandLineFile, _ := ioutil.TempFile("", "")

		suffix := HashSuf(config)
		username := fmt.Sprintf("user-%v", suffix)
		clustername := fmt.Sprintf("cluster-%v", suffix)

		for key, obj := range config.AuthInfos {
			config.AuthInfos[username] = obj
			delete(config.AuthInfos, key)
			break
		}
		for key, obj := range config.Clusters {
			config.Clusters[clustername] = obj
			delete(config.Clusters, key)
			break
		}
		for key, obj := range config.Contexts {
			obj.AuthInfo = username
			obj.Cluster = clustername
			config.Contexts[name] = obj
			delete(config.Contexts, key)
			break
		}
		configType := clientcmdapi.Config{
			AuthInfos: config.AuthInfos,
			Clusters:  config.Clusters,
			Contexts:  config.Contexts,
		}
		_ = clientcmd.WriteToFile(configType, commandLineFile.Name())
		loop = append(loop, commandLineFile.Name())
		mc.command.Printf("Context Add: %s \n", name)
		defaultName = name
	}
	loadingRules := &clientcmd.ClientConfigLoadingRules{
		Precedence: loop,
	}
	mergedConfig, err := loadingRules.Load()
	if mergedConfig != nil {
		mergedConfig.CurrentContext = defaultName
	}
	json, err := runtime.Encode(clientcmdlatest.Codec, mergedConfig)
	if err != nil {
		Error.Printf("Unexpected error: %v", err)
	}
	output, err := yaml.JSONToYAML(json)
	if err != nil {
		Error.Printf("Unexpected error: %v", err)
	}

	for _, name := range loop {
		defer os.Remove(name)
	}

	err = WriteConfig(output)
	if err != nil {
		Error.Println(err.Error())
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
