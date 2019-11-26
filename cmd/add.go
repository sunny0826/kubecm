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
	"github.com/spf13/cobra"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	"os"
	syaml "sigs.k8s.io/yaml"
	"strings"
)

var file string
var name string
var cover bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Merge configuration file with $HOME/.kube/config",
	Example: `
# Merge example.yaml with $HOME/.kube/config
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with $HOME/.kube/config
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileExists(file) {
			err := configCheck(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cover, _ = cmd.Flags().GetBool("cover")
			config, err := getAddConfig(file)
			if err != nil {
				fmt.Println(err)
			}
			output := merge2Master(config)
			err = WriteConfig(output)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Printf("%s file does not exist", file)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&file, "file", "f", "", "Path to merge kubeconfig files")
	addCmd.Flags().StringVarP(&name, "name", "n", "", "The name of contexts. if this field is null,it will be named with file name.")
	addCmd.Flags().BoolP("cover", "c", false, "Overwrite the original kubeconfig file")
	addCmd.MarkFlagRequired("file")
}

func getAddConfig(kubeconfig string) (*clientcmdapi.Config, error) {

	config, err := LoadClientConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	if len(config.AuthInfos) != 1 {
		fmt.Println("Only support add 1 context.")
		os.Exit(-1)
	}

	name := getName()
	err = nameCheck(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
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

	return config, nil
}

func nameCheck(name string) error {
	c, err := LoadClientConfig(cfgFile)
	if err != nil {
		return err
	}
	for key, _ := range c.Contexts {
		if key == name {
			return fmt.Errorf("The name: %s already exists, please replace it.", name)
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getName() string {
	if name == "" {
		n := strings.Split(file, "/")
		result := strings.Split(n[len(n)-1], ".")
		return result[0]
	} else {
		return name
	}
}

func configCheck(kubeconfigPath string) error {
	_, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	return nil
}

func LoadClientConfig(kubeconfig string) (*clientcmdapi.Config, error) {
	b, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.Load(b)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func merge2Master(config *clientcmdapi.Config) []byte {
	commandLineFile, _ := ioutil.TempFile("", "")
	defer os.Remove(commandLineFile.Name())
	configType := clientcmdapi.Config{
		AuthInfos: config.AuthInfos,
		Clusters:  config.Clusters,
		Contexts:  config.Contexts,
	}
	_ = clientcmd.WriteToFile(configType, commandLineFile.Name())
	loadingRules := &clientcmd.ClientConfigLoadingRules{
		Precedence: []string{cfgFile, commandLineFile.Name()},
	}

	mergedConfig, err := loadingRules.Load()

	json, err := runtime.Encode(clientcmdlatest.Codec, mergedConfig)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
	output, err := syaml.JSONToYAML(json)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}

	return output
}

func WriteConfig(config []byte) error {
	if cover {
		err := ioutil.WriteFile(cfgFile, config, 0777)
		if err != nil {
			return err
		}
	} else {
		err := ioutil.WriteFile("./config.yaml", config, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
